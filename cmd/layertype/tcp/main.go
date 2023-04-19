// Package main is the example of *layers.TCP
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/devlights/gomy/chans"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	appLog = log.New(os.Stderr, "", 0)
)

func main() {
	if err := run(); err != nil {
		appLog.Panic(err)
	}
}

func run() error {
	const (
		device      = "eth0"
		filter      = "tcp and port 443"
		snapshotLen = int32(1600)
		promiscuous = false
		timeout     = pcap.BlockForever
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
		dataSource    gopacket.PacketDataSource = handle
		decoder       gopacket.Decoder          = handle.LinkType()
		packetSource  *gopacket.PacketSource    = gopacket.NewPacketSource(dataSource, decoder)
		packetCh      <-chan gopacket.Packet    = packetSource.Packets()
		first3Packets                           = chans.Take(doneCh, packetCh, 3)
	)
	appLog.Println("START")

LOOP:
	for {
		select {
		case <-sigCh:
			close(doneCh)
			break LOOP
		case p, ok := <-first3Packets:
			if !ok {
				break LOOP
			}

			see(p)
		}
	}

	return nil
}

func see(p gopacket.Packet) {
	tcpLayer := p.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return
	}

	tcp := tcpLayer.(*layers.TCP)

	appLog.Printf("[Src Port       ] %v", tcp.SrcPort)
	appLog.Printf("[Dst Port       ] %v", tcp.DstPort)
	appLog.Printf("[Seq Number     ] %v", tcp.Seq)
	appLog.Printf("[Ack Number     ] %v", tcp.Ack)
	appLog.Printf("[Window Size    ] %v", tcp.Window)
	appLog.Printf("[TCP Flags - SYN] %v", tcp.SYN)
	appLog.Printf("[TCP Flags - ACK] %v", tcp.ACK)
	appLog.Printf("[TCP Flags - PSH] %v", tcp.PSH)
	appLog.Printf("[TCP Flags - RST] %v", tcp.RST)
	appLog.Printf("[TCP Flags - FIN] %v", tcp.FIN)
	appLog.Printf("[Checksum       ] %v", tcp.Checksum)
	appLog.Printf("[Urgent Pointer ] %v", tcp.Urgent)
	appLog.Println("----------------")
}
