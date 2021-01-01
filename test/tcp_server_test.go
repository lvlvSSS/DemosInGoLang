package test

import (
	"net"
	"testing"
	"time"
)

func TestResolveAddr(t *testing.T) {
	addr, _ := net.ResolveTCPAddr("", "127.0.0.1:50000")
	t.Logf("%v", addr.IP)

	var dur time.Duration
	t.Log(dur)
}
