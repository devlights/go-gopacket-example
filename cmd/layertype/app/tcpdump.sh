#!/usr/bin/env bash

tcpdump -i lo -w example.pcap 'tcp and port 22222' 1>/dev/null 2>&1 & 
echo -n $! > tcpdump.pid
exit 0