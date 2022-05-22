package cop

import "fmt"

var (
	// half-connection syn map
	hmap map[SocketId]*Socket
	// full-connection syn map
	fmap map[SocketId]*Socket
	// closing-connection fin map
	cmap map[SocketId]*Socket
)


func init() {
	hmap = make(map[SocketId]*Socket)
	fmap = make(map[SocketId]*Socket)
	cmap = make(map[SocketId]*Socket)
}



func push(sock *Socket) {
	switch sock.Status {
	case SYN_SENT:
		hmap[sock.sid] = sock
	case SYN_RCVD:
		fmap[sock.sid] = sock
	case FIN_WAIT_1:
		cmap[sock.sid] = sock
	}
}

func notClosed(sock *Socket) bool {
	return hmap[sock.sid] != nil || fmap[sock.sid] != nil || cmap[sock.sid] != nil
}

func pushHalfList(sock *Socket) error {
	if sock.Status != SYN_SENT && sock.Status != SYN_RCVD {
		return fmt.Errorf("Invalid socket status %v", sock)
	}
	if notClosed(sock) {
		return fmt.Errorf("Socket %v already exists", sock)
	}

	hmap[sock.sid] = sock
	return nil
}

func pushFullList(sock *Socket) error {
	return nil
}

func pushClosingList(sock *Socket) error {
	return nil
}
