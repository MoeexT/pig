package main

import (
	"flag"
	"fmt"
	"net"

	"pig/lib/cop"
	"pig/util"
	"pig/util/log"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)


var (
	si       string // src ip
	di       string // dst ip
	sp       uint   // src port
	dp       uint   // dst port
	a        bool   // all packets
	v        bool   // verbose
	protocol string // which protocol
	dog      *log.Logger
)

func init() {
	flag.StringVar(&di, "di", "", "dst ip")
	flag.StringVar(&si, "si", "", "src ip")
	flag.UintVar(&sp, "sp", 0, "src port")
	flag.UintVar(&dp, "dp", 0, "dst port")
	flag.BoolVar(&a, "a", false, "all packets")
	flag.BoolVar(&v, "v", false, "verbose mode: show headers")
	flag.StringVar(&protocol, "p", "tcp", "protocol: tcp/icmp")
	dog = log.Dog
}

func main() {
	flag.Parse()

	switch protocol {
	case "tcp":
		catchTcp()
	case "icmp":
		catchIcmp()
	default:
		flag.PrintDefaults()
	}
}

func catchTcp() {
	// srcAddr, _ := net.ResolveIPAddr("ip4", si)

	// catch all ip packets
	conn, _ := net.ListenIP("ip4:tcp", nil)
	ipConn, _ := ipv4.NewRawConn(conn)
	for {
		buf := make([]byte, 1500)
		ipHdr, payload, _, _ := ipConn.ReadFrom(buf)

		go func() {
			// output all ip headers
			if a {
				dog.Debug(ipHdr.String())
			}

			// filter ip
			if (si != "" && si != ipHdr.Src.String()) ||
				(di != "" && di != ipHdr.Dst.String()) {
				return
			}

			// parse tcp header
			hdr, data := cop.ParseHeader(payload)
			if hdr == nil {
				return
			}

			// filter port
			if (sp != 0 && sp != uint(hdr.SrcPort)) ||
				(dp != 0 && dp != uint(hdr.DstPort)) {
				return
			}
			if v {
				dog.Debug(hdr.String())
			}

			dog.Trace(fmt.Sprintf("%s:%d -> %s:%d", ipHdr.Src.String(), hdr.SrcPort, ipHdr.Dst.String(), hdr.DstPort))
			dog.Trace(fmt.Sprintf("%#v", ipHdr))
			dog.Info(util.BytesString(data, 20))
		}()
	}
}

func catchIcmp() {
	conn, _ := net.ListenIP("ip4:icmp", nil)
	for {
		buf := make([]byte, 1024)
		n, addr, _ := conn.ReadFrom(buf)
		go func() {
			msg, _ := icmp.ParseMessage(1, buf[0:n])
			dog.Info(n, addr, msg.Type, msg.Code, msg.Checksum)
		}()
	}
}
