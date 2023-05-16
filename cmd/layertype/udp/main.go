// Package main is the example of *layers.UDP
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

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
		device      = "lo"
		filter      = "udp and port 22222"
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
		dataSource    gopacket.PacketDataSource = handle
		decoder       gopacket.Decoder          = handle.LinkType()
		packetSource  *gopacket.PacketSource    = gopacket.NewPacketSource(dataSource, decoder)
		packetCh      <-chan gopacket.Packet    = packetSource.Packets()
		first1Packets                           = chans.Take(doneCh, packetCh, 1)
	)
	appLog.Println("START")

LOOP:
	for {
		select {
		case <-sigCh:
			close(doneCh)
			break LOOP
		case p, ok := <-first1Packets:
			if !ok {
				break LOOP
			}

			see(p)
		}
	}

	return nil
}

func see(p gopacket.Packet) {
	udpLayer := p.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return
	}

	udp := udpLayer.(*layers.UDP)

	appLog.Printf("[Src Port       ] %v", udp.SrcPort)
	appLog.Printf("[Dst Port       ] %v", udp.DstPort)
	appLog.Printf("[Length         ] %v", udp.Length)
	appLog.Printf("[Payload        ] %v", udp.Payload)
	appLog.Printf("[Payload(decode)] %v", string(udp.Payload))
	appLog.Printf("[Checksum       ] %v", udp.Checksum)
	appLog.Println("----------------")
}
