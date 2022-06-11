package mip

import (
    "errors"
    "net"
    "strconv"
    "strings"
    "syscall"

    "pig/lib/eth"
    "pig/util/logger"

    "golang.org/x/net/ipv4"
)

type IPv4Addr = [4]byte
type IPv6Addr = [16]byte

var (
    fd        int
    addr      *syscall.SockaddrInet4
    addrCache map[uint64]*syscall.SockaddrInet4
    dog       *logger.Logger
)

func init() {
    fd, _ = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
    dog = logger.Dog
    addrCache = make(map[uint64]*syscall.SockaddrInet4)
}

func NewHeader(src, dst net.IP, dl int) *ipv4.Header {
    return &ipv4.Header{
        Version:  4,
        Len:      20,
        TOS:      0,
        TotalLen: 20 + dl,
        ID:       0,
        FragOff:  0,
        Flags:    0b10,
        TTL:      64,
        Protocol: 6,
        Checksum: 0,
        Src:      src,
        Dst:      dst,
    }
}

// turn an ipv4 address string to IPv4Addr
func ParseIPv4(addr string) (*IPv4Addr, error) {
    ss := strings.Split(addr, ".")
    var ip4ad = new(IPv4Addr)
    if len(ss) != 4 {
        return nil, errors.New("invalid ip format")
    }

    for i := 0; i < 4; i++ {
        frag, err := strconv.ParseInt(ss[i], 10, 16)

        if frag < 0 || frag > 255 || err != nil {
            return nil, errors.New("invalid ip format")
        }
        ip4ad[i] = byte(frag)
    }
    return ip4ad, nil
}

func SendTo(ip net.IP, port uint16, content []byte) error {
    id := (uint64(ip[0]) << 40) + (uint64(ip[1]) << 32) + (uint64(ip[2]) << 24) + (uint64(ip[3]) << 16) + uint64(port)
    addr = addrCache[id]
    if addr == nil {
        addr = &syscall.SockaddrInet4{
            Port: int(port),
            Addr: [4]byte{ip[0], ip[1], ip[2], ip[3]},
        }
        addrCache[id] = addr
        dog.Trace("new addr:", addr)
    }
    dog.Tracef("send to %v:%v, content: %v", ip, port, content)
    return syscall.Sendto(fd, content, 0, addr)
}

func htons(i uint16) uint16 {
    return (i&0xff)<<8 | i>>8
}

// receive all ipv4 packets from the network then yield data
func BeginReceive() (ch chan []byte, err error) {
    fdr, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_IP)))
    dog.Trace("created socket with fd:", fdr)
    if err != nil {
        return nil, err
    }

    ch = make(chan []byte, 256)

    ifi, err := net.InterfaceByName("eth0")
    dog.Info("interface:", ifi)

    go func() {
        for {
            buf := make([]byte, 1518)
            n, _, err := syscall.Recvfrom(fdr, buf, 0)
            if err != nil || n < 14 {
                dog.Errorf("read %d bytes from ethernet with error: %v", n, err)
                continue
            }

            // if llsa, ok := addr.(*syscall.SockaddrLinklayer); ok {
            // 	inter, err := net.InterfaceByIndex(llsa.Ifindex)
            // 	if err != nil {
            // 		dog.Error(os.Stderr, "interface from ifindex: %s", err.Error())
            // 	}
            // 	dog.Info(inter.Name + ": ")
            // }

            eHeader := eth.ParseHeader(buf[:14])

            if eHeader.EtherType != eth.IPv4 {
                dog.Warn("not ipv4 packet: " + eHeader.String())
                ch <- buf[14:n]
            } else {
                ch <- buf[14:n]
            }
        }
    }()

    return ch, nil
}

func Close() {
    syscall.Close(fd)
}
