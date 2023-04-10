// Package main is the example of go-packet with BPF (Berkeley Packet Filter) filter
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/google/gopacket"
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
		pcapfile = "example.pcap"
		filter   = "icmp or icmp6"
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
	// Apply capture filter
	//
	// # filter examples:
	//   - ip src 192.168.1.1
	//   - ip dst 192.168.1.2
	//   - ip host 192.168.1.1 and ip host 192.168.1.2
	//   - tcp port 80
	//   - udp port 53
	//   - icmp or icmp6
	//   - ether src aa:bb:cc:dd:ee:ff
	//   - ether dst aa:bb:cc:dd:ee:ff
	//   - vlan 100
	//   - ip host 192.168.1.1 and tcp port 80
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

			appLog.Println(p)
		}
	}

	return nil
}
