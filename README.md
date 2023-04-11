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

### nc (netcat)

```sh
$ sudo apt install netcat
```

### tcpdump

```sh
$ sudo apt install tcpdump
```

### [go-task](https://taskfile.dev/)

```sh
$ go install github.com/go-task/task/v3/cmd/task@latest
```

## How to run

```sh
$ task --list
task: Available tasks for this project:
* bpffilter:         Run pcap.OpenOffline() with BPF Filter
* default:           default (print all ifs)
* fmtvet:            go fmt and go vet
* openlive:          Run pcap.OpenLive() example
* openoffline:       Run pcap.OpenOffline() example
* packet:            See gopacket.Packet structure info


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
```

## REFERENCES

- [gopacket](https://pkg.go.dev/github.com/google/gopacket@v1.1.19)
- [Sniffing packets in Go](https://medium.com/a-bit-off/sniffing-network-go-6753cae91d3f)
- [gopacketでpcapを読み込む](https://mrtc0.hateblo.jp/entry/2016/03/19/232252)
- [ncコマンドでサービスの接続疎通確認](https://qiita.com/chenglin/items/70f06e146db19de5a659)
