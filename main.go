package main

import (
	"bufio"
	"os"
	"pig/lib/cop"
	"pig/util/log"
)

/*

https://toutiao.io/posts/wj9ori/preview

*/

const (
	EmptyString = ""
	HELP_INFO = `custom optimization protocol
	:conn to connect.
	:disc to disconnect.
	:help to show help.
	:stat to show socket's status.
	:sl[t|d|i|w|e|f] to set log level.
	:exit,:quit to exit program.`
)

var (
	input string
	dog   *log.Logger
)

func main() {
	dog = log.Dog
	dog.Level = log.Info
	dog.Info(HELP_INFO)

	reader := bufio.NewReader(os.Stdin)

	chRcv := make(chan []byte, 128)

	// sender
	sock := cop.NewSocket([]byte{127, 0, 0, 1}, []byte{127, 0, 0, 1}, 626, 627)


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
			dog.Level = log.Trace
		case ":sld":
			dog.Level = log.Debug
		case ":sli":
			dog.Level = log.Info
		case ":slw":
			dog.Level = log.Warn
		case ":sle":
			dog.Level = log.Error
		case ":slf":
			dog.Level = log.Fatal
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
