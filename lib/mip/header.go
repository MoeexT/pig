package mip

import (
	"golang.org/x/net/ipv4"
)


func Marshal(h *ipv4.Header, content []byte) (bytes []byte) {
	/* bytes = make([]byte, h.TotalLen)
	bytes[0] = (byte(h.Version) << 4) + (byte(h.Len) >> 2)
	bytes[1] = byte(h.TOS)
	bytes[2] = byte((h.TotalLen & 0xff00) >> 8)
	bytes[3] = byte(h.TotalLen & 0xff)
	bytes[4] = byte((h.ID & 0xff00) >> 8)
	bytes[5] = byte(h.ID & 0xff)
	bytes[6] = byte((int(h.Flags) << 5) + ((h.FragOff & 0x1f00) >> 8))
	bytes[7] = byte(h.FragOff & 0xff)
	bytes[8] = byte(h.TTL)
	bytes[9] = byte(h.Protocol)
	bytes[10] = 0
	bytes[11] = 0
	bytes[12] = h.Src[0]
	bytes[13] = h.Src[1]
	bytes[14] = h.Src[2]
	bytes[15] = h.Src[3]
	bytes[16] = h.Dst[0]
	bytes[17] = h.Dst[1]
	bytes[18] = h.Dst[2]
	bytes[19] = h.Dst[3] */
	bytes = append(bytes, content...)
	h.Checksum = int(Checksum(bytes))
	bytes, _ = h.Marshal()
	bytes = append(bytes, content...)
	return
}

// cop checksum calculation
func Checksum(bytes []byte) uint16 {
	bts := bytes[:]
	bts[16] = 0
	bts[17] = 0

	sum := 0
	for n := 1; n < len(bytes)-1; n += 2 {
		sum += (int(bytes[n]) << 8) + int(bytes[n+1])
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += sum >> 16
	return uint16(^sum)
}
