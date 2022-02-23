package main

import (
	"flag"
	"fmt"
	"net"

	"pig/lib/cop"

	"golang.org/x/net/ipv4"
)

var (
	s uint   // src port
	d uint   // dst port
	a bool   // all packets
	i string // ip
)

func init() {
	flag.UintVar(&s, "s", 0, "src port")
	flag.UintVar(&d, "d", 0, "dst port")
	flag.BoolVar(&a, "a", false, "all packets")
	flag.StringVar(&i, "i", "127.0.0.1", "ip")
}

func main() {
	flag.Parse()
	if !a && s == 0 && d == 0 {
		flag.PrintDefaults()
		return
	}

	addr, _ := net.ResolveIPAddr("ip4", i)
	fmt.Println("addr: ", addr.IP, addr.Zone)
	conn, _ := net.ListenIP("ip4:tcp", addr)
	ipConn, _ := ipv4.NewRawConn(conn)
	for {
		buf := make([]byte, 1500)
		ipHdr, payload, _, _ := ipConn.ReadFrom(buf)
		// hdr, payload, controlMessage, _ := ipConn.ReadFrom(buf)
		// fmt.Println(hdr, len(payload), controlMessage)
		hdr, data := cop.ParseHeader(payload)

		if hdr == nil {
			continue
		}

		if a {
			fmt.Println(ipHdr)
			fmt.Println(buf[:20])
		} else if hdr.SrcPort == uint16(s) || hdr.DstPort == uint16(d) {
			fmt.Println(string(data))
		}
	}
}
