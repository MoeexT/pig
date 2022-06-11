package main

import (
    "flag"
    "net"
    "os"

    "pig/lib/cop"
    "pig/util/logger"

    "golang.org/x/term"
)

var (
    si  string // source ip, my ip
    di  string // dst ip
    sp  uint   // source port, my port
    dp  uint   // dst port
    dog *logger.Logger
)

const (
    HELP_INFO = "\r\n--- acho --- a program just can echo" +
        "\r\n\t't'/'d'/'i'/'w'/'e'/'f' to set log level." +
        "\r\n\t's' to show socket's status." +
        "\r\n\t'h' to show help." +
        "\r\n\t'q' to quit program.\r"
)

func init() {
    flag.StringVar(&si, "si", "", "source ip")
    flag.StringVar(&di, "di", "", "dst ip")
    flag.UintVar(&sp, "sp", 0, "source port")
    flag.UintVar(&dp, "dp", 0, "dst port")
    dog = logger.Dog
    dog.Level = logger.Trace
}

func main() {
    flag.Parse()
    if si == "" || di == "" {
        flag.PrintDefaults()
        return
    }
    dog.Info(HELP_INFO)
    dog.Infof("source: %v:%v", si, sp)
    dog.Infof("dstnation: %v:%v", di, dp)

    sa, err := net.ResolveIPAddr("ip4", si)
    if err != nil {
        dog.Error("ResolveIPAddr error: %v", err)
    }
    da, err := net.ResolveIPAddr("ip4", di)
    if err != nil {
        dog.Error("ResolveIPAddr error: %v", err)
    }
    sock := cop.NewSocket(sa.IP, da.IP, uint16(sp), uint16(dp))
    rch := make(chan []byte, 128)
    go sock.BeginReceive(rch)
    go func() {
        for {
            data := <-rch
            dog.Info(string(data))
            go sock.Send(data)
        }
    }()

    // set terminal raw mode
    oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
    dog.Debugf("oldState: %v", oldState)
    if err != nil {
        dog.Errorf("MakeRaw error: %v", err)
        return
    }
    defer term.Restore(int(os.Stdin.Fd()), oldState)
    // term.Restore(int(os.Stdin.Fd()), oldState)

    c := make([]byte, 1)
    go func() {
        for {
            os.Stdin.Read(c)

            switch c[0] {
            case 't':
                dog.Level = logger.Trace
            case 'd':
                dog.Level = logger.Debug
            case 'i':
                dog.Level = logger.Info
            case 'w':
                dog.Level = logger.Warn
            case 'e':
                dog.Level = logger.Error
            case 'f':
                dog.Level = logger.Fatal
            case 's':
                dog.Info(cop.SOCK_STAT_DESC[sock.Status])
            case 'h':
                dog.Info(HELP_INFO)
            case 'q':
                dog.Warn("sure to quit?")
                os.Stdin.Read(c)
                if c[0] == 'y' {
                    err = term.Restore(int(os.Stdin.Fd()), oldState)
                    dog.Fatalf("quit with error: %v", err)
                    os.Exit(0)
                }
            }
        }
    }()
    select {}
}
