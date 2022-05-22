package main

import (
	"bufio"
	"flag"
	"net"
	"os"

	"pig/lib/cop"
	"pig/util/logger"
)

/*

https://toutiao.io/posts/wj9ori/preview

*/

const (
	EmptyString = ""
	HELP_INFO   = `custom optimization protocol
	:conn to connect.
	:disc to disconnect.
	:help to show help.
	:stat to show socket's status.
	:sl[t|d|i|w|e|f] to set log level.
	:exit,:quit to exit program.`
)

var (
	si       string // src ip
	di       string // dst ip
	sp       uint   // src port
	dp       uint   // dst port
	a        bool   // all packets
	v        bool   // verbose
	protocol string // which protocol
	input    string
	dog      *logger.Logger
)

func init() {
	flag.StringVar(&di, "di", "127.0.0.1", "dst ip")
	flag.StringVar(&si, "si", "127.0.0.1", "src ip")
	flag.UintVar(&sp, "sp", 626, "src port")
	flag.UintVar(&dp, "dp", 627, "dst port")
	flag.BoolVar(&a, "a", false, "all packets")
	flag.BoolVar(&v, "v", false, "verbose mode: show headers")
	flag.StringVar(&protocol, "p", "tcp", "protocol: tcp/icmp")
	dog = logger.Dog
}

func main() {
	flag.Parse()
	flag.PrintDefaults()
	dog = logger.Dog
	dog.Level = logger.Info
	dog.Info(HELP_INFO)

	reader := bufio.NewReader(os.Stdin)

	chRcv := make(chan []byte, 128)

	// sender
	sock := cop.NewSocket(net.ParseIP(si), net.ParseIP(di), uint16(sp), uint16(dp))

	// receive message
	go sock.BeginReceive(chRcv)
	go func() {
		for {
			dog.Info(string(<-chRcv))
		}
	}()

	// send message
	for {
		// get input from console stdin
		input, _ = reader.ReadString('\n')
		if len(input) == 1 {
			continue
		}

		input = input[:len(input)-1]

		switch input {
		case ":conn":
			err := sock.Connect()
			if err != nil {
				dog.Errorf("connect error: %v", err)
			}
			continue
		case ":slt":
			dog.Level = logger.Trace
		case ":sld":
			dog.Level = logger.Debug
		case ":sli":
			dog.Level = logger.Info
		case ":slw":
			dog.Level = logger.Warn
		case ":sle":
			dog.Level = logger.Error
		case ":slf":
			dog.Level = logger.Fatal
		case ":stat":
			dog.Info(cop.SOCK_STAT_DESC[sock.Status])
		case ":help":
			dog.Info(HELP_INFO)
		case ":disc":
			err := sock.Close()
			if err != nil {
				dog.Errorf("close error: %v", err)
			}
			continue
		case ":exit", ":quit":
			dog.Warn("Bye")
			return
		default:
			go sock.Send([]byte(input))
		}
		input = EmptyString
	}
}
