package cop

import "fmt"

type HeaderFlag struct {
}

type Header struct {
	// src port
	SrcPort uint16
	// dst port
	DstPort uint16
	// sequence number
	Seq uint32
	// acknowledge number
	Ack uint32
	// header len
	Offset uint8
	// flags such as ACK, SYN, FIN
	Flags uint16
	// window size
	WinSize uint16
	// checksum
	Checksum uint16
	// urgent pointer
	UrgPtr uint16
	// tcp options
	Options []byte
}

func (h *Header) String() string {
	if h == nil {
		return "nil"
	}
	return fmt.Sprintf("COP.Header{src: %d, dst: %d, seq: %d, ack: %d, offset: %d, flags: %09b, win-size: %d, check-sum: %04x, urgent: %d, options: %v}",
		h.SrcPort, h.DstPort, h.Seq, h.Ack, h.Offset, h.Flags, h.WinSize, h.Checksum, h.UrgPtr, h.Options)
}

func (h *Header) MarshalSelf() (bytes []byte) {
	bytes = make([]byte, 20)

	bytes[0] = byte(h.SrcPort >> 8)
	bytes[1] = byte(h.SrcPort)
	bytes[2] = byte(h.DstPort >> 8)
	bytes[3] = byte(h.DstPort)
	bytes[4] = byte(h.Seq >> 24)
	bytes[5] = byte(h.Seq >> 16)
	bytes[6] = byte(h.Seq >> 8)
	bytes[7] = byte(h.Seq)
	bytes[8] = byte(h.Ack >> 24)
	bytes[9] = byte(h.Ack >> 16)
	bytes[10] = byte(h.Ack >> 8)
	bytes[11] = byte(h.Ack)
	bytes[12] = byte((h.Offset << 2) + uint8((h.Flags&0x100)>>8))
	bytes[13] = byte(h.Flags)
	bytes[14] = byte(h.WinSize >> 8)
	bytes[15] = byte(h.WinSize)
	bytes[16] = byte(h.Checksum >> 8)
	bytes[17] = byte(h.Checksum)

	return
}

func (h *Header) Marshal(content []byte) (bytes []byte) {
	bytes = make([]byte, 20)

	bytes[0] = byte(h.SrcPort >> 8)
	bytes[1] = byte(h.SrcPort)
	bytes[2] = byte(h.DstPort >> 8)
	bytes[3] = byte(h.DstPort)
	bytes[4] = byte(h.Seq >> 24)
	bytes[5] = byte(h.Seq >> 16)
	bytes[6] = byte(h.Seq >> 8)
	bytes[7] = byte(h.Seq)
	bytes[8] = byte(h.Ack >> 24)
	bytes[9] = byte(h.Ack >> 16)
	bytes[10] = byte(h.Ack >> 8)
	bytes[11] = byte(h.Ack)
	bytes[12] = byte((h.Offset << 2) + uint8((h.Flags&0x100)>>8))
	bytes[13] = byte(h.Flags)
	bytes[14] = byte(h.WinSize >> 8)
	bytes[15] = byte(h.WinSize)
	// checksum and Urgent pointer are both 0

	bytes = append(bytes, content...)

	cs := Checksum(bytes)
	bytes[16] = byte((cs & 0xff00) >> 8)
	bytes[17] = byte(cs & 0xff)

	return
}

// parse bytes-header to cop header
func ParseHeader(bytes []byte) (header *Header, err error) {
	if bytes == nil || len(bytes) < 20 {
		return nil, fmt.Errorf("invalid header bytes")
	}

	header = new(Header)
	// validate header
	header.Checksum = (uint16(bytes[16]) << 8) + uint16(bytes[17])
	// if header.Checksum != Checksum(bytes) {
	// 	fmt.Printf("%04x, %04x\n", header.Checksum, Checksum(bytes))
	// 	return nil
	// }

	header.SrcPort = (uint16(bytes[0]) << 8) + uint16(bytes[1])
	header.DstPort = (uint16(bytes[2]) << 8) + uint16(bytes[3])
	header.Seq = (uint32(bytes[4]) << 24) + (uint32(bytes[5]) << 16) + (uint32(bytes[6]) << 8) + uint32(bytes[7])
	header.Ack = (uint32(bytes[8]) << 24) + (uint32(bytes[9]) << 16) + (uint32(bytes[10]) << 8) + uint32(bytes[11])
	header.Offset = (bytes[12] & 0xf0) >> 2
	header.Flags = ((uint16(bytes[12]) & 1) << 8) + uint16(bytes[13])
	header.WinSize = (uint16(bytes[14]) << 8) + uint16(bytes[15])
	header.UrgPtr = (uint16(bytes[18]) << 8) + uint16(bytes[19])

	// tcp header options, if has
	if header.Offset > 20 {
		header.Options = bytes[20:header.Offset]
	}

	return
}

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
