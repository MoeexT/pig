package mip_test

import (
	"pig/lib/mip"
	"testing"
)

func TestIPAddressParser(t *testing.T) {
	addr, err := mip.ParseIPv4("-255.0.0.1")
	t.Log(err)
	t.Log(addr)
}
