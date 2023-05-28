// Package main is the example of HTTP using go-packet.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
		device      = "lo"
		filter      = "tcp and port 12345"
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
				if !errors.Is(err, NoPayload) {
					return err
				}
			}
		}
	}

	return nil
}

var (
	NoPayload = errors.New("no payload")
)

func see(p gopacket.Packet) error {
	tcpLayer := p.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return nil
	}

	tcp := tcpLayer.(*layers.TCP)

	payload := tcp.LayerPayload()
	if len(payload) == 0 {
		return NoPayload
	}

	// ペイロードを文字列に変換してHTTPとしての解析を試みる
	payloadStr := string(payload)
	reader := bufio.NewReader(strings.NewReader(payloadStr))

	if isResponse(&payloadStr) {
		//
		// HTTPレスポンス
		//
		resp, err := http.ReadResponse(reader, nil)
		if err != nil {
			return err
		}

		appLog.Println("HTTP Status Code:", resp.StatusCode)
		appLog.Println("HTTP Protocol:", resp.Proto)
		appLog.Println("HTTP Headers:")
		for name, values := range resp.Header {
			appLog.Printf("  %s: %s\n", name, strings.Join(values, ", "))
		}
		appLog.Println()
	} else if isRequest(&payloadStr) {
		//
		// HTTPリクエスト
		//
		req, err := http.ReadRequest(reader)
		if err == nil {
			appLog.Println("HTTP Method:", req.Method)
			appLog.Println("HTTP URL:", req.URL)
			appLog.Println("HTTP Protocol:", req.Proto)
			appLog.Println("HTTP Headers:")
			for name, values := range req.Header {
				appLog.Printf("  %s: %s\n", name, strings.Join(values, ", "))
			}
			appLog.Println()
		}
	}

	return nil
}

func isResponse(s *string) bool {
	return strings.HasPrefix(*s, "HTTP")
}

var (
	methods = []string{
		"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "CONNECT", "PATCH",
	}
)

func isRequest(s *string) bool {

	target := *s
	for _, method := range methods {
		if strings.HasPrefix(target, method) {
			return true
		}
	}

	return false
}
