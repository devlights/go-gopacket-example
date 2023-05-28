// Package main is the example of HTTP using go-packet.
package main

import (
	"bufio"
	"fmt"
	"io"
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
	var (
		tcpLayer gopacket.Layer
	)

	tcpLayer = p.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return nil
	}

	//
	// ペイロードを取得
	//
	var (
		tcp     = tcpLayer.(*layers.TCP)
		payload = tcp.LayerPayload()
	)

	if len(payload) == 0 {
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
		appLog.Println("------------------------------------")
		return nil
	}

	//
	// ペイロードを文字列に変換してHTTPとしての解析を試みる
	//
	var (
		payloadStr = string(payload)
		reader     = bufio.NewReader(strings.NewReader(payloadStr))
	)
	defer func() { appLog.Println("------------------------------------") }()

	if isRequest(&payloadStr) {
		//
		// HTTPリクエスト
		//
		var (
			req *http.Request
			err error
		)

		req, err = http.ReadRequest(reader)
		if err != nil {
			return err
		}

		appLog.Println("HTTP Method:", req.Method)
		appLog.Println("HTTP URL:", req.URL)
		appLog.Println("HTTP Protocol:", req.Proto)

		appLog.Println("HTTP Headers:")
		for name, values := range req.Header {
			appLog.Printf("  %s: %s\n", name, strings.Join(values, ", "))
		}
	}

	if isResponse(&payloadStr) {
		//
		// HTTPレスポンス
		//
		var (
			resp *http.Response
			body []byte
			err  error
		)

		resp, err = http.ReadResponse(reader, nil)
		if err != nil {
			return err
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		appLog.Println("HTTP Status Code:", resp.StatusCode)
		appLog.Println("HTTP Protocol:", resp.Proto)

		appLog.Println("HTTP Headers:")
		for name, values := range resp.Header {
			appLog.Printf("  %s: %s\n", name, strings.Join(values, ", "))
		}

		appLog.Println("BODY:", string(body))
	}

	return nil
}

func isResponse(s *string) bool {
	return strings.HasPrefix(*s, "HTTP")
}

func isRequest(s *string) bool {
	var (
		methods = []string{
			"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "CONNECT", "PATCH",
		}
		target = *s
	)

	for _, method := range methods {
		if strings.HasPrefix(target, method) {
			return true
		}
	}

	return false
}
