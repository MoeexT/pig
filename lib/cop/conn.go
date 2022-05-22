package cop

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"pig/lib/mip"
	"pig/util"
	"pig/util/log"

	"golang.org/x/net/ipv4"
)

const (
	// socket status
	SYN_SENT = 1 << iota
	SYN_RCVD
	ESTABLISHED
	FIN_WAIT_1
	FIN_WAIT_2
	CLOSE_WAIT
	LAST_ACK
	TIME_WAIT
	CLOSED = 0

	// time wait span: 1 sec
	TIME_WAIT_TIME = 1

	// tcp flags
	SYN1_FLAGS = 0b00010
	SYN2_FLAGS = 0b10010
	// SYN3_FLAGS = 0b10000 == ACK
	ACK_FLAGS = 0b10000
	FIN_FLAGS = 0b10001
	// FIN2_FLAGS = 0b10000 == ACK
	// FIN3_FLAGS = 0b10001 == FIN
	// FIN4_FLAGS = 0b10000 == ACK

)

var (
	// syn channels
	// synchs map[SocketId]chan error
	logger         *log.Logger
	SOCK_STAT_DESC map[uint8]string
)

func init() {
	logger = log.Dog

	SOCK_STAT_DESC = map[uint8]string{
		SYN_SENT:    "SYN_SENT",
		SYN_RCVD:    "SYN_RCVD",
		ESTABLISHED: "ESTABLISHED",
		FIN_WAIT_1:  "FIN_WAIT_1",
		FIN_WAIT_2:  "FIN_WAIT_2",
		CLOSE_WAIT:  "CLOSE_WAIT",
		LAST_ACK:    "LAST_ACK",
		TIME_WAIT:   "TIME_WAIT",
		CLOSED:      "CLOSED",
	}
}

type SocketId [12]byte
type Socket struct {
	sid SocketId // socket id
	seq uint32   // sequence number
	nxt *Socket  // next socket, used for linked-list-management

	Status  uint8 // connection status
	SrcIP   []byte
	DstIP   []byte
	SrcPort uint16
	DstPort uint16
}

func NewSocket(srcIP, dstIP []byte, srcPort, dstPort uint16) *Socket {
	sock := &Socket{
		SrcIP:   srcIP,
		DstIP:   dstIP,
		SrcPort: srcPort,
		DstPort: dstPort,
		Status:  CLOSED,
	}
	sock.sid = sock.id()
	sock.seq = rand.Uint32()
	return sock
}

// connect to remote end
func (sock *Socket) Connect() (err error) {
	// 1. send SYN
	return sock.synchronize()
	// 2. get response of this SYN
	// let it go
	// process this step in `replySyn`

	// 3. send ACK
	// process in `acknowledge`
}

func (sock *Socket) newHeader(flag uint16) *Header {
	return &Header{
		SrcPort: sock.SrcPort,
		DstPort: sock.DstPort,
		Seq:     sock.seq,
		Ack:     0x655f4,
		Offset:  20,
		Flags:   flag,
		WinSize: 0xff,
	}
}

func (sock *Socket) Send(content []byte) (err error) {
	if sock.Status > LAST_ACK {
		logger.Error("Socket %v is not established", sock)
		return fmt.Errorf("socket is not ESTABLISHED")
	}
	copHeader := Header{
		SrcPort: sock.SrcPort,
		DstPort: sock.DstPort,
		Seq:     sock.seq,
		Ack:     0x655f4,
		Offset:  20,
		Flags:   0b10000,
		WinSize: 0xff,
	}
	copPayload := copHeader.Marshal(content)

	ipHdr := mip.NewHeader(sock.SrcIP, sock.DstIP, len(copPayload))
	ipPayload := mip.Marshal(ipHdr, copPayload)

	err = mip.SendTo(sock.DstIP, sock.DstPort, ipPayload)

	// validation
	phdr, _ := ipv4.ParseHeader(ipPayload)
	logger.Tracef("ip header: %v", phdr)
	logger.Tracef("ip send(%d): %v", len(ipPayload), util.BytesString(ipPayload, len(ipPayload)))
	return
}

func (sock *Socket) BeginReceive(ch chan []byte) (err error) {
	addr := &net.IPAddr{
		IP: sock.SrcIP,
	}
	conn, err := net.ListenIP("ip4:tcp", addr)
	if err != nil {
		logger.Errorf("ListenIP error: %v", err)
		return
	}

	ipConn, err := ipv4.NewRawConn(conn)
	if err != nil {
		logger.Errorf("NewRawConn error: %v", err)
		return
	}

	for {
		buf := make([]byte, 1500)
		ipHeader, payload, _, err := ipConn.ReadFrom(buf)
		if err != nil {
			logger.Errorf("Read IP packet error: %v", err)
			ch <- nil
		}

		copHeader, err := ParseHeader(payload)
		if err != nil {
			dog.Error(err)
			continue
		}
		
		// filter dst port
		if copHeader.DstPort != sock.SrcPort {
			// logger.Tracef("dst port[%d] not match: %d", sock.SrcPort, copHeader.DstPort)
			continue
		}

		data := payload[copHeader.Offset:]
		// data is nil, as if this control packet
		if len(data) == 0 {
			// ack, syn, fin
			switch copHeader.Flags & 0b10011 {
			// first handshake
			case SYN1_FLAGS:
				logger.Debug("receive SYN")
				go sock.replySyn1()
			// second handshake
			case SYN2_FLAGS:
				logger.Debug("receive SYN-ACK")
				go sock.replySyn2()
			// third handshake, second handwave, third handwave
			case ACK_FLAGS:
				logger.Debug("receive ACK")
				go sock.replyAck()
			// fin
			case FIN_FLAGS:
				logger.Debug("receive FIN")
				go sock.replyFin()
			}
		} else {
			// data packet
			logger.Trace("IP Header", ipHeader)
			logger.Trace("cop Header", copHeader)
			ch <- data
			sock.acknowledge()
		}
	}
}

// first handshake, send SYN
func (sock *Socket) synchronize() (err error) {
	logger.Debug("Send SYN")

	// check if the connection is established or establishing
	// or, check the status
	if hmap[sock.sid] != nil || fmap[sock.sid] != nil || cmap[sock.sid] != nil {
		return fmt.Errorf("Socket %v already exists", sock)
	}
	if sock.Status != CLOSED {
		return fmt.Errorf("Socket %v already exists", sock)
	}

	err = sock.sendControl(SYN1_FLAGS)

	// set status
	sock.Status = SYN_SENT

	// add socket to half-connection map
	hmap[sock.sid] = sock

	return
}

// second shakehand, send SYN-ACK
// 收到第一次握手报文，设置状态SYN_RCVD
// 将这个连接放到半连接队列里
func (sock *Socket) replySyn1() (err error) {

	// if hmap[sock.sid] != nil || fmap[sock.sid] != nil || cmap[sock.sid] != nil {
	// 	return fmt.Errorf("Socket %v already exists", sock)
	// }
	if sock.Status != CLOSED {
		logger.Error("Socket %v already exists", sock)
		return fmt.Errorf("Socket %v already exists", sock)
	}

	logger.Debug("Reply SYN1")
	// set status
	sock.Status = SYN_RCVD
	// hmap[sock.sid] = sock
	return sock.sendControl(SYN2_FLAGS)
}

// 收到SYN-ACK报文，设置状态ESTABLISHED
func (sock *Socket) replySyn2() (err error) {
	sock.Status = ESTABLISHED
	// fmap[sock.sid] = hmap[sock.sid]
	// delete(hmap, sock.sid)
	logger.Info("Established")
	return sock.sendControl(ACK_FLAGS)
}

func (sock *Socket) replyAck() {
	switch sock.Status {
	case SYN_RCVD:
		// 3rd handshake
		logger.Info("Established")
		sock.Status = ESTABLISHED
	case FIN_WAIT_1:
		logger.Debug("Reply ACK FIN_WAIT_1")
		// 2nd handwave
		sock.Status = FIN_WAIT_2
	case LAST_ACK:
		logger.Info("Closed")
		// 4th handwave
		sock.Status = CLOSED
	}
}

// send ACK
func (sock *Socket) acknowledge() (err error) {
	logger.Debug("Send ACK")
	return sock.sendControl(ACK_FLAGS)
}

func (sock *Socket) replyFin() (err error) {
	logger.Debug("Reply FIN")

	switch sock.Status {
	case ESTABLISHED:
		logger.Debug("ESTABLISHED to close_wait")
		// 1st handwave
		sock.Status = CLOSE_WAIT
		// 2nd handwave
		sock.sendControl(ACK_FLAGS)
		// 这里能继续发送数据
		sock.Send([]byte("Bye"))
		// 3rd handwave
		sock.Status = LAST_ACK
		sock.sendControl(FIN_FLAGS)
	case FIN_WAIT_2:
		logger.Debug("fin_wait_2 to time_wait")
		sock.Status = TIME_WAIT
		// 4th handwave
		sock.sendControl(ACK_FLAGS)
		go sock.timeWait()
	}
	return

	// if fmap[sock.sid] != nil {
	// 	// first FIN
	// 	// set status, move socket to closing-connection map
	// 	sock.Status = CLOSE_WAIT
	// 	cmap[sock.sid] = fmap[sock.sid]
	// 	delete(fmap, sock.sid)
	// } else if cmap[sock.sid] != nil {
	// 	// second FIN
	// 	// set status, remove socket from closing-connection map in seconds
	// 	sock.Status = TIME_WAIT
	// 	go sock.timeWait()
	// } else { // unexpected circumstance
	// 	return fmt.Errorf("Socket %v has disconnected", sock)
	// }
}

// finish connection
func (sock *Socket) finish() (err error) {
	logger.Info("Finish")

	if sock.Status != ESTABLISHED {
		return fmt.Errorf("Socket %v has disconnected", sock)
	}

	err = sock.sendControl(FIN_FLAGS)
	sock.Status = FIN_WAIT_1
	return

	// if fmap[sock.sid] != nil {
	// } else {
	// 	return fmt.Errorf("Socket %v has disconnected", sock)
	// }
	// cmap[sock.sid] = sock
}

func (sock *Socket) sendControl(flags uint16) error {
	finCopHdr := sock.newHeader(flags)
	finCopPld := finCopHdr.Marshal(nil)
	finIpHdr := mip.NewHeader(sock.SrcIP, sock.DstIP, 0)
	finIpPld := mip.Marshal(finIpHdr, finCopPld)
	return mip.SendTo(sock.DstIP, sock.DstPort, finIpPld)
}

func (sock *Socket) timeWait() {
	time.Sleep(time.Second * TIME_WAIT_TIME)
	// delete(cmap, sock.sid)
	logger.Info("Closed")
	sock.Status = CLOSED
	mip.Close()
}

// generate socket id
func (sock *Socket) id() SocketId {
	bi := new(SocketId)
	bi[0] = sock.SrcIP[0]
	bi[1] = sock.SrcIP[1]
	bi[2] = sock.SrcIP[2]
	bi[3] = sock.SrcIP[3]
	bi[4] = byte(sock.SrcPort >> 8)
	bi[5] = byte(sock.SrcPort & 0xff)
	bi[6] = sock.DstIP[0]
	bi[7] = sock.DstIP[1]
	bi[8] = sock.DstIP[2]
	bi[9] = sock.DstIP[3]
	bi[10] = byte(sock.DstPort >> 8)
	bi[11] = byte(sock.DstPort & 0xff)
	return *bi
}

func (sock *Socket) Close() error {
	return sock.finish()
}
