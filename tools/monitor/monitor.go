package main

import (
	"flag"
	"fmt"
	"net"

	"pig/lib/cop"
	"pig/lib/mip"
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
	flag.StringVar(&protocol, "p", "ip", "protocol: tcp/icmp")
	dog = log.Dog
}

func main() {
	flag.Parse()

	switch protocol {
	case "ip":
		catchIP()
	case "tcp":
		catchTcpWithChan()
		// catchTcp()
	case "icmp":
		catchIcmp()
	default:
		flag.PrintDefaults()
	}
}

func catchIP() {
	ch, err := mip.BeginReceive()
	if err != nil {
		panic(err)
	}

	for frame := range ch {
		ipHdr, err := ipv4.ParseHeader(frame)
		if err != nil {
			dog.Error("read from ipv4 failed:", err)
		}
		dog.Debug("ipv4 header:", ipHdr)
		copHeader, err := cop.ParseHeader(frame[ipHdr.Len:])

		if err == nil {
			dog.Infof("COP header: %v",copHeader)
			segment := frame[ipHdr.Len+int(copHeader.Offset):]

			// filter port
			if (sp != 0 && sp != uint(copHeader.SrcPort)) ||
				(dp != 0 && dp != uint(copHeader.DstPort)) {
				return
			}
			if v {
				dog.Debug(copHeader.String())
			}

			dog.Trace(fmt.Sprintf("%s:%d -> %s:%d", ipHdr.Src.String(),
				copHeader.SrcPort, ipHdr.Dst.String(), copHeader.DstPort))
			dog.Info(util.BytesString(segment, 20))
		} else {
			dog.Error(err)
		}
	}
}

func catchTcpWithChan() {
	ch, err := cop.BeginRead()
	if err != nil {
		dog.Error("BeginReceive error:", err)
		return
	}

	for packet := range ch {
		hdr, err := cop.ParseHeader(packet)
		if err != nil {
			data := packet[hdr.Offset:]
			dog.Info("TCP header:", hdr, len(data))
		}
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

			// parse tcp header
			hdr, err := cop.ParseHeader(payload)
			if err != nil {
				dog.Error(err)
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
			dog.Info(util.BytesString(payload[hdr.Offset:], 20))
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
