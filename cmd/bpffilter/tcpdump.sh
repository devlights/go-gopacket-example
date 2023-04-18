#!/usr/bin/env bash

tcpdump -t -n -i lo -w example.pcap &
echo $! > tcpdump.pid
exit 0
