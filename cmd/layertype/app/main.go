// Package main is the example of gopacket.Packet.ApplicationLayer()
package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/devlights/gomy/logops"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	appLog, errLog, _ = logops.Default.Logger(false)
)

func main() {
	if err := run(); err != nil {
		errLog.Panic(err)
	}
}

func run() error {
	const (
		pcapfile = "example.pcap"
		filter   = ""
	)

	defer func() { appLog.Println("DONE") }()

	// --------------------------------------
	// Open capture handle
	// --------------------------------------
	var (
		handle *pcap.Handle
		err    error
	)

	handle, err = pcap.OpenOffline(pcapfile)
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

			see(p)
		}
	}

	return nil
}

func see(p gopacket.Packet) {
	//
	// 以下、どちらも同じことをしている
	//

	// gopacket.Packet.ApplicationLayer() で、アプリケーション層のペイロードを取得
	appLayer := p.ApplicationLayer()
	if appLayer == nil {
		return
	}

	appLog.Printf("[ApplicatonLayer][Payload ] %v bytes", len(appLayer.Payload()))
	appLog.Printf("[ApplicatonLayer][Contents] %v", string(appLayer.Payload()))

	// TCPレイヤーからペイロードを取得
	tcpLayer := p.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return
	}

	tcp := tcpLayer.(*layers.TCP)

	appLog.Printf("[TCP Layer      ][Payload ] %v bytes", len(tcp.Payload))
	appLog.Printf("[TCP Layer      ][Contents] %v", string(tcp.Payload))
	appLog.Println("----------------")
}
