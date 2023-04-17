// Package main is the example of *layers.ARP
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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
		filter      = ""
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
		sigCh = make(chan os.Signal, 1)
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
			break LOOP
		case p, ok := <-packetCh:
			if !ok {
				break LOOP
			}

			// Display only the first items
			if see(p) {
				break LOOP
			}
		}
	}

	return nil
}

func see(p gopacket.Packet) bool {
	arpLayer := p.Layer(layers.LayerTypeARP)
	if arpLayer == nil {
		return false
	}

	arp := arpLayer.(*layers.ARP)

	appLog.Printf("[Operation    ] %v", arp.Operation)
	appLog.Printf("[Src Hw Addr  ] %v", arp.SourceHwAddress)
	appLog.Printf("[Src Prot Addr] %v", arp.SourceProtAddress)
	appLog.Printf("[Dst Hw Addr  ] %v", arp.DstHwAddress)
	appLog.Printf("[Dst Prot Addr] %v", arp.DstProtAddress)

	return true
}
