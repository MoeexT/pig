package cop_test

import (
	"net"
	"pig/lib/cop"
	"pig/util"
	"testing"

	"golang.org/x/net/ipv4"
)

func TestCopHeader(t *testing.T) {
	header := cop.Header{
		SrcPort:  43318,
		DstPort:  10541,
		Seq:      248246597,
		Ack:      1707584022,
		Offset:   20,
		Flags:    0x18,
		WinSize:  647,
		Checksum: 0xef05,
		UrgPtr:   0,
	}
	hb := header.MarshalSelf()
	t.Log(util.BytesString(hb, 20))
}

func TestIPHeader(t *testing.T) {
	ipHdr := ipv4.Header{
		Version:  4,
		Len:      20,
		TotalLen: 353,
		ID:       0x3988,
		TOS:      0,
		TTL:      44,
		Flags:    0x40,
		FragOff:  0,
		Protocol: 6,
		Checksum: 0x90ba,
		Src:      net.IPv4(119, 28, 43, 34),
		Dst:      net.IPv4(192, 168, 32, 110),
	}
	ipBytes, _ := ipHdr.Marshal()
	t.Log(util.BytesString(ipBytes, 20))
}
