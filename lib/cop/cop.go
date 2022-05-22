package cop

import (
	"net"

	"pig/util/logger"

	"golang.org/x/net/ipv4"
)

type IpPacket struct {
	Header  *ipv4.Header
	Payload []byte
	Control *ipv4.ControlMessage
}

var (
	dog *logger.Logger
)

func init() {
	dog = logger.Dog
}

func ConnectTo() {

}

func BeginRead() (ch chan []byte, err error) {
	conn, err := net.ListenIP("ip4:tcp", nil)
	if err != nil {
		return nil, err
	}

	ch = make(chan []byte)

	go func() {
		for {
			buf := make([]byte, 1500)
			n, addr, err := conn.ReadFromIP(buf)

			if err == nil {
				dog.Trace("read from tcp:", addr, n)
				ch <- buf[:n]
			} else {
				dog.Errorf("parse ipv4(%dB) header failed: %v", n, err)
			}
		}
	}()

	return ch, nil
}

func ReceiveWithNet() (ch chan *IpPacket, err error) {
	conn, err := net.ListenIP("ip4:tcp", nil)
	if err != nil {
		dog.Errorf("listen ipv4 failed: %v", err)
		return nil, nil
	}
	icConn, err := ipv4.NewRawConn(conn)
	if err != nil {
		dog.Errorf("new ipv4 raw conn failed: %v", err)
		return nil, nil
	}

	go func() {
		buf := make([]byte, 1500)
		for {
			ipHdr, pld, cm, err := icConn.ReadFrom(buf)
			if err == nil {
				ch <- &IpPacket{
					Header:  ipHdr,
					Payload: pld,
					Control: cm,
				}
			} else {
				dog.Errorf("read from ipv4 failed: %v", err)
			}
		}
	}()

	return ch, nil
}
