package test

import (
    "fmt"
    "runtime"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestJust(t *testing.T) {
    x := 0x2333
    t.Log(fmt.Sprintf("0x%X", byte(x)))
    t.Log(fmt.Sprintf("0x%X", byte(x>>8)))
    assert.Equal(t, byte(x>>8), byte((x&0xff00)>>8))
    assert.NotEqual(t, byte(x&0xff00), byte((x&0xff00)>>8))

    y := 0x123456789abcdef
    assert.Equal(t, byte(0xde), byte(y>>4))
    assert.Equal(t, byte(0xcd), byte(y>>8))
    assert.Equal(t, byte(0xbc), byte(y>>12))
    assert.Equal(t, byte(0xab), byte(y>>16))
    assert.Equal(t, byte(0x9a), byte(y>>20))
    assert.Equal(t, byte(0x89), byte(y>>24))
    assert.Equal(t, byte(0x78), byte(y>>28))
    assert.Equal(t, byte(0x67), byte(y>>32))
    assert.Equal(t, byte(0x56), byte(y>>36))
    assert.Equal(t, byte(0x45), byte(y>>40))
    assert.Equal(t, byte(0x34), byte(y>>44))
    assert.Equal(t, byte(0x23), byte(y>>48))
    assert.Equal(t, byte(0x12), byte(y>>52))
    assert.Equal(t, byte(0x01), byte(y>>56))
    assert.Equal(t, byte(0x00), byte(y>>60))
}

func TestCaller(t *testing.T) {
    pc, _, _, ok := runtime.Caller(0)
    fmt.Println(ok, runtime.FuncForPC(pc).Name())
    pc, _, _, ok = runtime.Caller(1)
    fmt.Println(ok, runtime.FuncForPC(pc).Name())
    pc, _, _, ok = runtime.Caller(2)
    fmt.Println(ok, runtime.FuncForPC(pc).Name())
    pc, _, _, ok = runtime.Caller(3)
    fmt.Println(ok, runtime.FuncForPC(pc).Name())
    pc, _, _, ok = runtime.Caller(4)
    fmt.Println(ok, runtime.FuncForPC(pc).Name())
    pc, _, _, ok = runtime.Caller(5)
    fmt.Println(ok, runtime.FuncForPC(pc).Name())

}

func TestIota(t *testing.T) {
    const (
        SYN_SENT = 1 << iota
        SYN_RCVD
        ESTABLISHED
        FIN_WAIT_1
        FIN_WAIT_2
        CLOSE_WAIT
        LAST_ACK
        TIME_WAIT
        CLOSED = 0
    )

    fmt.Println("CLOSED: ", CLOSED)
    fmt.Println("SYN_SENT: ", SYN_SENT)
    fmt.Println("SYN_RCVD: ", SYN_RCVD)
    fmt.Println("ESTABLISHED: ", ESTABLISHED)
    fmt.Println("FIN_WAIT_1: ", FIN_WAIT_1)
    fmt.Println("FIN_WAIT_2: ", FIN_WAIT_2)
    fmt.Println("CLOSE_WAIT: ", CLOSE_WAIT)
    fmt.Println("LAST_ACK: ", LAST_ACK)
    fmt.Println("TIME_WAIT: ", TIME_WAIT)
}
