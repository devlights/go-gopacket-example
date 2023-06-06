// Package main is the example of DNS using go-packet.
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
	"github.com/miekg/dns"
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
		filter      = "udp and port 53"
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
	udpLayer := p.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return nil
	}

	//
	// ペイロードを取得
	//
	var (
		udp     = udpLayer.(*layers.UDP)
		payload = udp.LayerPayload()
	)

	appLog.Printf("[Src Port       ] %v", udp.SrcPort)
	appLog.Printf("[Dst Port       ] %v", udp.DstPort)
	appLog.Printf("[Length         ] %v", udp.Length)
	appLog.Printf("[Checksum       ] %v", udp.Checksum)

	if len(payload) == 0 {
		return nil
	}

	defer func() { appLog.Println("------------------------------------") }()

	//
	// DNSペイロードとして解釈する
	//
	var (
		msg = new(dns.Msg)
		err error
	)

	err = msg.Unpack(payload)
	if err != nil {
		return err
	}

	var (
		hasQuestion = len(msg.Question) > 0
		hasAnswer   = len(msg.Answer) > 0
	)

	// DNS Question
	if hasQuestion {
		appLog.Println("[DNS Questions]")
		for _, q := range msg.Question {
			appLog.Printf("\t%v\n", q.String())
		}
	}

	// DNS Answer
	if hasAnswer {
		appLog.Println("[DNS Answers]")
		for _, a := range msg.Answer {
			appLog.Printf("\t%v\n", a.String())
		}
	}

	return nil
}
