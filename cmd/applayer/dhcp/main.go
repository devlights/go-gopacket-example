// Package main is the example of DHCP using go-packet.
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	appLog = log.New(os.Stderr, "", 0)
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	const (
		device      = "eth0"
		filter      = "udp and (port 67 or port 68)"
		snapshotLen = int32(1600)
		promiscuous = false
		timeout     = 1 * time.Second
	)

	defer func() { appLog.Println("DONE") }()

	// --------------------------------------
	// Open capture handle
	// --------------------------------------
	var (
		handle *pcap.Handle
		err    error
	)

	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		return fmt.Errorf("error open handle: %w", err)
	}
	defer handle.Close()

	// --------------------------------------
	// Apply capture filter (optional)
	// --------------------------------------
	if filter != "" {
		err = handle.SetBPFFilter(filter)
		if err != nil {
			return fmt.Errorf("error apply filter: %w", err)
		}
	}

	// --------------------------------------
	// Set signal handler
	// --------------------------------------
	var (
		doneCh = make(chan struct{})
		sigCh  = make(chan os.Signal, 1)
	)

	signal.Notify(sigCh, os.Interrupt)

	// --------------------------------------
	// Make packet source and display.
	// --------------------------------------
	var (
		dataSource   gopacket.PacketDataSource = handle
		decoder      gopacket.Decoder          = handle.LinkType()
		packetSource *gopacket.PacketSource    = gopacket.NewPacketSource(dataSource, decoder)
		packetCh     <-chan gopacket.Packet    = packetSource.Packets()
	)
	appLog.Println("START")

LOOP:
	for {
		select {
		case <-sigCh:
			close(doneCh)
			break LOOP
		case p, ok := <-packetCh:
			if !ok {
				break LOOP
			}

			if err = see(p); err != nil {
				return err
			}
		}
	}

	return nil
}

func see(p gopacket.Packet) error {
	//
	// レイヤーを確認
	//
	dhcpLayer := p.Layer(layers.LayerTypeDHCPv4)
	if dhcpLayer == nil {
		return nil
	}

	//
	// DHCPパケットに変換
	//
	dhcp, ok := dhcpLayer.(*layers.DHCPv4)
	if !ok {
		return fmt.Errorf("fail: convert *layers.DHCPv4")
	}

	appLog.Printf("[Operation          ] %v\n", dhcp.Operation)
	appLog.Printf("[Hardware Type      ] %v\n", dhcp.HardwareType)
	appLog.Printf("[Hardware Length    ] %v\n", dhcp.HardwareLen)
	appLog.Printf("[Hardware Options   ] %v\n", dhcp.HardwareOpts)
	appLog.Printf("[DHCP Xid           ] %v\n", dhcp.Xid)
	appLog.Printf("[DHCP Secs          ] %v\n", dhcp.Secs)
	appLog.Printf("[DHCP Flags         ] %v\n", dhcp.Flags)
	appLog.Printf("[DHCP Client IP     ] %v\n", dhcp.ClientIP)
	appLog.Printf("[DHCP Your Client IP] %v\n", dhcp.YourClientIP)
	appLog.Printf("[DHCP Next Server IP] %v\n", dhcp.NextServerIP)
	appLog.Printf("[DHCP Relay Agent IP] %v\n", dhcp.RelayAgentIP)
	appLog.Printf("[DHCP Client HW Addr] %v\n", dhcp.ClientHWAddr)
	appLog.Printf("[DHCP Server Name   ] %v\n", dhcp.ServerName)
	appLog.Printf("[DHCP File          ] %v\n", dhcp.File)

	for _, dhcpOption := range dhcp.Options {
		appLog.Printf("\t[DHCP Option     ] %v\n", dhcpOption.Type)
		appLog.Printf("\t[DHCP Option Data] %v\n", dhcpOption.Data)
	}

	defer func() { appLog.Println("------------------------------------") }()

	return nil
}
