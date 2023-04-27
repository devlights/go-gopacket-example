# go-gopacket-example

Packet Capture with gopacket example by golang.

![Go Version](https://img.shields.io/badge/go-1.20-blue.svg)

The sources in this repository only work on Linux.

## Environments

```sh
$ lsb_release -a
No LSB modules are available.
Distributor ID: Ubuntu
Description:    Ubuntu 22.04.2 LTS
Release:        22.04
Codename:       jammy


$ go version
go version go1.20.3 linux/amd64
```

## Requirements

### libpcap

```sh
$ sudo apt install libpcap-dev
```

### nc (netcat) (optional)

```sh
$ sudo apt install netcat
```

### tcpdump (optional)

```sh
$ sudo apt install tcpdump
```

### arp-scan (optional)

```sh
$ sudo apt install arp-scan
```

### [go-task](https://taskfile.dev/)

```sh
$ go install github.com/go-task/task/v3/cmd/task@latest
```

## How to run

```sh
$ task --list
task: Available tasks for this project:
* bpffilter:                Run pcap.OpenOffline() with BPF Filter
* default:                  default (print all ifs)
* fmtvet:                   go fmt and go vet
* layertype-app:            See gopacket.Packet.ApplicationLayer() info
* layertype-arp:            See *layers.ARP info
* layertype-ethernet:       See *layers.Ethernet info
* layertype-icmpv4:         See *layers.ICMPv4 info
* layertype-ipv4:           See *layers.IPv4 info
* layertype-tcp:            See *layers.TCP info
* layertype-udp:            See *layers.UDP info
* openlive:                 Run pcap.OpenLive() example
* openoffline:              Run pcap.OpenOffline() example
* packet:                   See *pcap.Packet structure info


$ task openlive
task: [openlive] go build
task: [openlive] sudo ./openlive


[Packet capture will be displayed.]


$ task openoffline
task: [openoffline] sudo timeout 3s tcpdump -i eth0 -w example.pcap 'tcp'
tcpdump: listening on eth0, link-type EN10MB (Ethernet), capture size 262144 bytes
61 packets captured
77 packets received by filter
0 packets dropped by kernel
task: [openoffline] go build
task: [openoffline] sudo ./openoffline

[Packet capture will be displayed.]



$ task bpffilter
task: [bpffilter] go build
task: [bpffilter] sudo bash ./tcpdump.sh
task: [bpffilter] bash ./ping.sh
PING localhost(localhost (::1)) 56 data bytes
64 bytes from localhost (::1): icmp_seq=1 ttl=64 time=0.018 ms
tcpdump: listening on lo, link-type EN10MB (Ethernet), snapshot length 262144 bytes
64 bytes from localhost (::1): icmp_seq=2 ttl=64 time=0.031 ms
64 bytes from localhost (::1): icmp_seq=3 ttl=64 time=0.046 ms
task: [bpffilter] sudo bash ./kill.sh
33 packets captured
82 packets received by filter
0 packets dropped by kernel
task: [bpffilter] sleep 1
task: [bpffilter] sudo ./bpffilter
START

[Packet capture will be displayed.]

DONE



$ task packet
task: [packet] go build
task: [packet] sudo bash ./app.sh
task: [packet] sleep 1
task: [packet] sudo bash ./server.sh
task: [packet] sudo bash ./client.sh
helloworldtask: [packet] sleep 3
------------------------------
[Capture Length] 74
[Src           ] 127.0.0.1
[Dst           ] 127.0.0.1
[Protocol      ] TCP
[SRC PORT      ] 45400
[DST PORT      ] 22222(easyengine)
[TCP FLAGS     ]
>>> SYN=true
>>> ACK=false
>>> PSH=false
>>> RST=false
>>> FIN=false
------------------------------
------------------------------
[Capture Length] 74
[Src           ] 127.0.0.1
[Dst           ] 127.0.0.1
[Protocol      ] TCP
[SRC PORT      ] 22222(easyengine)
[DST PORT      ] 45400
[TCP FLAGS     ]
>>> SYN=true
>>> ACK=true
>>> PSH=false
>>> RST=false
>>> FIN=false
------------------------------
------------------------------
[Capture Length] 66
[Src           ] 127.0.0.1
[Dst           ] 127.0.0.1
[Protocol      ] TCP
[SRC PORT      ] 45400
[DST PORT      ] 22222(easyengine)
[TCP FLAGS     ]
>>> SYN=false
>>> ACK=true
>>> PSH=false
>>> RST=false
>>> FIN=false
------------------------------
------------------------------
[Capture Length] 76
[Src           ] 127.0.0.1
[Dst           ] 127.0.0.1
[Protocol      ] TCP
[SRC PORT      ] 45400
[DST PORT      ] 22222(easyengine)
[TCP FLAGS     ]
>>> SYN=false
>>> ACK=true
>>> PSH=true
>>> RST=false
>>> FIN=false
[Payload       ] [104 101 108 108 111 119 111 114 108 100]
------------------------------
------------------------------
[Capture Length] 66
[Src           ] 127.0.0.1
[Dst           ] 127.0.0.1
[Protocol      ] TCP
[SRC PORT      ] 22222(easyengine)
[DST PORT      ] 45400
[TCP FLAGS     ]
>>> SYN=false
>>> ACK=true
>>> PSH=false
>>> RST=false
>>> FIN=false
------------------------------
------------------------------
[Capture Length] 66
[Src           ] 127.0.0.1
[Dst           ] 127.0.0.1
[Protocol      ] TCP
[SRC PORT      ] 45400
[DST PORT      ] 22222(easyengine)
[TCP FLAGS     ]
>>> SYN=false
>>> ACK=true
>>> PSH=false
>>> RST=false
>>> FIN=true
------------------------------
------------------------------
[Capture Length] 66
[Src           ] 127.0.0.1
[Dst           ] 127.0.0.1
[Protocol      ] TCP
[SRC PORT      ] 22222(easyengine)
[DST PORT      ] 45400
[TCP FLAGS     ]
>>> SYN=false
>>> ACK=true
>>> PSH=false
>>> RST=false
>>> FIN=true
------------------------------
------------------------------
[Capture Length] 66
[Src           ] 127.0.0.1
[Dst           ] 127.0.0.1
[Protocol      ] TCP
[SRC PORT      ] 45400
[DST PORT      ] 22222(easyengine)
[TCP FLAGS     ]
>>> SYN=false
>>> ACK=true
>>> PSH=false
>>> RST=false
>>> FIN=false
------------------------------
task: [packet] sudo bash ./kill.sh



$ task layertype-ethernet
task: [layertype-ethernet] go build
task: [layertype-ethernet] sudo ./ethernet
START
[Src MAC      ] 16:xx:42:44:xx:cd
[Dst MAC      ] 7e:bf:24:xx:3e:90
[Ethernet type] IPv4
[Src MAC      ] 7e:bf:24:xx:3e:90
[Dst MAC      ] 16:xx:42:44:2e:cd
[Ethernet type] IPv4
[Src MAC      ] 16:xx:42:44:2e:cd
[Dst MAC      ] 7e:bf:24:xx:3e:90
[Ethernet type] IPv4
[Src MAC      ] 7e:bf:24:xx:3e:90
[Dst MAC      ] 16:xx:42:44:2e:cd
[Ethernet type] IPv4
[Src MAC      ] 7e:bf:24:xx:3e:90
[Dst MAC      ] 16:xx:42:44:2e:cd
[Ethernet type] IPv4
[Src MAC      ] 7e:bf:24:xx:3e:90
[Dst MAC      ] 16:xx:42:44:2e:cd
[Ethernet type] IPv4
DONE


$ task layertype-arp
task: [layertype-arp] go build
task: [layertype-arp] sudo bash ./arp-scan.sh &
task: [layertype-arp] sudo ./arp
START
Interface: eth0, type: EN10MB, MAC: c2:88:65:43:bc:ed, IPv4: 10.0.5.2
Starting arp-scan 1.9.7 with 4 hosts (https://github.com/royhills/arp-scan)
10.0.5.1        42:6f:a6:72:06:80       (Unknown: locally administered)
[Operation    ] 1
[Src Hw Addr  ] [194 136 101 76 188 237]
[Src Prot Addr] [10 0 2 6]
[Dst Hw Addr  ] [0 0 0 0 0 0]
[Dst Prot Addr] [10 3 5 0]

1 packets received by filter, 0 packets dropped by kernel
Ending arp-scan 1.9.7: 4 hosts scanned in 1.443 seconds (2.77 hosts/sec). 1 responded
DONE


$ task layertype-ipv4
task: [layertype-ipv4] go build
task: [layertype-ipv4] sudo ./ipv4
START
[Version       ] 4
[IHL           ] 5 words -> 160 bits -> 20 bytes
[Length        ] 81
[Payload Length] 61
[TTL           ] 64
[Protocol      ] TCP
[Src IP        ] 10.0.5.2
[Dst IP        ] 192.168.39.75
DONE 


$ task layertype-icmpv4
task: [layertype-icmpv4] go build
task: [layertype-icmpv4] sudo bash ./ping.sh &
task: [layertype-icmpv4] sudo ./icmpv4
START
PING  (127.0.0.1) 56(84) bytes of data.
64 bytes from localhost (127.0.0.1): icmp_seq=1 ttl=64 time=0.018 ms
[Seq     ] 1
[Type    ] 8
[Code    ] 0
[Req/Rep ] ICMP Echo Request
[Checksum] 20183
[Seq     ] 1
[Type    ] 0
[Code    ] 0
[Req/Rep ] ICMP Echo Reply
[Checksum] 22231
64 bytes from localhost (127.0.0.1): icmp_seq=2 ttl=64 time=0.028 ms

---  ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1020ms
rtt min/avg/max/mdev = 0.018/0.023/0.028/0.005 ms
DONE


$ task layertype-tcp
task: [layertype-tcp] go build
task: [layertype-tcp] sudo ./tcp
START
[Src Port       ] 51190
[Dst Port       ] 443(https)
[Seq Number     ] 771501779
[Ack Number     ] 1030112796
[Window Size    ] 331
[TCP Flags - SYN] false
[TCP Flags - ACK] true
[TCP Flags - PSH] true
[TCP Flags - RST] false
[TCP Flags - FIN] false
[Checksum       ] 4529
[Urgent Pointer ] 0
----------------
[Src Port       ] 51190
[Dst Port       ] 443(https)
[Seq Number     ] 771501819
[Ack Number     ] 1030112796
[Window Size    ] 331
[TCP Flags - SYN] false
[TCP Flags - ACK] true
[TCP Flags - PSH] true
[TCP Flags - RST] false
[TCP Flags - FIN] false
[Checksum       ] 4686
[Urgent Pointer ] 0
----------------
[Src Port       ] 443(https)
[Dst Port       ] 51190
[Seq Number     ] 1030112796
[Ack Number     ] 771501819
[Window Size    ] 1962
[TCP Flags - SYN] false
[TCP Flags - ACK] true
[TCP Flags - PSH] false
[TCP Flags - RST] false
[TCP Flags - FIN] false
[Checksum       ] 10449
[Urgent Pointer ] 0
----------------
DONE


$ task layertype-udp
task: [layertype-udp] go build
task: [layertype-udp] sudo bash ./app.sh
task: [layertype-udp] sleep 1
START
task: [layertype-udp] sudo bash ./server.sh
task: [layertype-udp] sudo bash ./client.sh
task: [layertype-udp] sleep 3
[Src Port       ] 38037
[Dst Port       ] 22222
[Length         ] 19
[Payload        ] [104 101 108 108 111 119 111 114 108 100 10]
[Payload(decode)] helloworld
[Checksum       ] 65062
----------------
task: [layertype-udp] sudo bash ./kill.sh



$ task layertype-app
task: [layertype-app] go build
task: [layertype-app] sudo bash ./tcpdump.sh
task: [layertype-app] sudo bash ./server.sh
task: [layertype-app] sleep 1
task: [layertype-app] sudo bash ./client.sh
task: [layertype-app] sleep 3
task: [layertype-app] sudo bash ./kill.sh
task: [layertype-app] sleep 1
task: [layertype-app] sudo ./app
START
[ApplicatonLayer][Payload ] 10 bytes
[ApplicatonLayer][Contents] helloworld
[TCP Layer      ][Payload ] 10 bytes
[TCP Layer      ][Contents] helloworld
----------------
[ApplicatonLayer][Payload ] 6 bytes
[ApplicatonLayer][Contents] golang
[TCP Layer      ][Payload ] 6 bytes
[TCP Layer      ][Contents] golang
----------------
[ApplicatonLayer][Payload ] 9 bytes
[ApplicatonLayer][Contents] goroutine
[TCP Layer      ][Payload ] 9 bytes
[TCP Layer      ][Contents] goroutine
----------------
DONE
```

## REFERENCES

- [gopacket](https://pkg.go.dev/github.com/google/gopacket@v1.1.19)
- [Sniffing packets in Go](https://medium.com/a-bit-off/sniffing-network-go-6753cae91d3f)
- [gopacketでpcapを読み込む](https://mrtc0.hateblo.jp/entry/2016/03/19/232252)
- [ncコマンドでサービスの接続疎通確認](https://qiita.com/chenglin/items/70f06e146db19de5a659)
- [IPが分からないオンプレミスをコマンドラインから調べる。(arp-scan)](https://qiita.com/iganari/items/7be4681ecfa5cff76feb)
