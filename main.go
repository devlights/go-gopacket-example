package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket/pcap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var (
		ifs []pcap.Interface
		err error
	)

	ifs, err = pcap.FindAllDevs()
	if err != nil {
		return err
	}

	for _, aIf := range ifs {
		fmt.Printf("Name=%-20vDescription=%v\n", aIf.Name, aIf.Description)
	}

	return nil
}
