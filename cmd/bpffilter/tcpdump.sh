#!/usr/bin/env bash

tcpdump -i lo -w example.pcap &
echo $! > tcpdump.pid
exit 0
