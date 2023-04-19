// Package main is the example of gopacket.Packet
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
		device      = "lo"
		filter      = "tcp port 22222"
		snapshotLen = int32(128)
		promiscuous = false
	)

	// --------------------------------------
	// Open capture handle
	// --------------------------------------
	var (
		handle *pcap.Handle
		err    error
	)

	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, pcap.BlockForever)
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

LOOP:
	for {
		select {
		case <-sigCh:
			break LOOP
		case p, ok := <-packetCh:
			if !ok {
				break LOOP
			}

			see(p)
		}
	}

	return nil
}

func see(p gopacket.Packet) {
	// ダンプ出力（各レイヤー毎の詳細も見れる）
	//appLog.Printf("[Dump] %v", p.Dump())
	// データ (各レイヤー毎のフルパケットデータが見れる)
	//appLog.Printf("[Data] %v", p.Data())

	appLog.Println("------------------------------")
	{
		appLog.Printf("[Capture Length] %v", p.Metadata().CaptureLength)

		ipLayer := p.Layer(layers.LayerTypeIPv4)
		if ipLayer != nil {
			ipv4, ok := ipLayer.(*layers.IPv4)
			if ok {
				appLog.Printf("[Src           ] %v", ipv4.SrcIP)
				appLog.Printf("[Dst           ] %v", ipv4.DstIP)
				appLog.Printf("[Protocol      ] %v", ipv4.Protocol)
			}
		}

		tcpLayer := p.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			tcp, ok := tcpLayer.(*layers.TCP)
			if ok {
				appLog.Printf("[SRC PORT      ] %v", tcp.SrcPort)
				appLog.Printf("[DST PORT      ] %v", tcp.DstPort)
				appLog.Println("[TCP FLAGS     ]")
				appLog.Printf(">>> SYN=%v", tcp.SYN)
				appLog.Printf(">>> ACK=%v", tcp.ACK)
				appLog.Printf(">>> PSH=%v", tcp.PSH)
				appLog.Printf(">>> RST=%v", tcp.RST)
				appLog.Printf(">>> FIN=%v", tcp.FIN)
			}
		}

		appLayer := p.ApplicationLayer()
		if appLayer != nil {
			payload := appLayer.Payload()
			if payload != nil {
				appLog.Printf("[Payload       ] %v", payload)
			}
		}

	}
	appLog.Println("------------------------------")
}
