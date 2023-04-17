// Package main is the example of *layers.ICMPv4
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
		appLog.Panic(err)
	}
}

func run() error {
	const (
		device      = "lo"
		filter      = "icmp"
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
	for i := 0; ; {
		select {
		case <-sigCh:
			break LOOP
		case p, ok := <-packetCh:
			if !ok {
				break LOOP
			}

			// Display only the first 2 items
			if see(p) {
				i++
				if i >= 2 {
					break LOOP
				}
			}
		}
	}

	return nil
}

func see(p gopacket.Packet) bool {
	icmpLayer := p.Layer(layers.LayerTypeICMPv4)
	if icmpLayer == nil {
		return false
	}

	icmpv4 := icmpLayer.(*layers.ICMPv4)

	// Ping（ICMP エコーリクエスト）の Type は 8 で、Code は 0 となる
	// 宛先ホストからのエコー応答の Type は 0 で、Code も 0 となる
	appLog.Printf("[Seq     ] %v", icmpv4.Seq)
	appLog.Printf("[Type    ] %v", icmpv4.TypeCode.Type())
	appLog.Printf("[Code    ] %v", icmpv4.TypeCode.Code())
	appLog.Printf("[Req/Rep ] %v", getType(icmpv4))
	appLog.Printf("[Checksum] %v", icmpv4.Checksum)

	return true
}

func getType(icmp *layers.ICMPv4) string {
	t := icmp.TypeCode
	if t.Type() == 8 && t.Code() == 0 {
		return "ICMP Echo Request"
	}

	if t.Type() == 0 && t.Code() == 0 {
		return "ICMP Echo Reply"
	}

	return ""
}
