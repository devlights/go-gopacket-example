// Package main is the example of *layers.IPv4
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
		device      = "eth0"
		filter      = "tcp"
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
	ipv4Layer := p.Layer(layers.LayerTypeIPv4)
	if ipv4Layer == nil {
		return false
	}

	ipv4 := ipv4Layer.(*layers.IPv4)

	appLog.Printf("[Version       ] %v", ipv4.Version)

	// Internet Header Length (IHL)
	//   IPv4ヘッダの長さを32ビットワード単位で示す.
	//   IHLは通常、最小値の5（5 * 32ビット = 160ビット = 20バイト）であるが
	//   IPヘッダにオプションが含まれている場合、それより大きくなる.
	//
	//   IHLが5の場合、IPv4ヘッダは20バイトの長さを持つ.
	//   IHLが6の場合、IPv4ヘッダは24バイトの長さを持ち、そのうち4バイトがオプションフィールドに割り当てられる.
	//   IHLを使用して、IPヘッダの終わりとペイロード（データ）の開始を正確に判断することができる.
	appLog.Printf("[IHL           ] %v words -> %v bits -> %v bytes", ipv4.IHL, ipv4.IHL*32, ipv4.IHL*32/8)

	// Lengthは IPv4 パケット全体での長さを表す
	appLog.Printf("[Length        ] %v", ipv4.Length)
	appLog.Printf("[Payload Length] %v", len(ipv4.Payload))

	appLog.Printf("[TTL           ] %v", ipv4.TTL)
	appLog.Printf("[Protocol      ] %v", ipv4.Protocol)
	appLog.Printf("[Src IP        ] %v", ipv4.SrcIP)
	appLog.Printf("[Dst IP        ] %v", ipv4.DstIP)

	return true
}
