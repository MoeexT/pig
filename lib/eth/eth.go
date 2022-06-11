package eth

import (
    "encoding/binary"
    "fmt"

    "pig/util"
)

type Header struct {
    DstMac    []byte
    SrcMac    []byte
    EtherType uint16
}

const (
    IPv4 = 0x0800
)

func (hdr *Header) String() string {
    if hdr == nil {
        return "nil"
    }
    return fmt.Sprintf("dst: %v, src: %v, ether-type: %v",
        util.BytesString(hdr.DstMac, len(hdr.DstMac)),
        util.BytesString(hdr.SrcMac, len(hdr.SrcMac)),
        hdr.EtherType)
}

func ParseHeader(bytes []byte) *Header {
    return &Header{
        DstMac:    bytes[0:6],
        SrcMac:    bytes[6:12],
        EtherType: binary.BigEndian.Uint16(bytes[12:14]),
    }
}
