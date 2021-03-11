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

type changeRef struct{
	number int
}

func TestChangeRef(t *testing.T){
	a := &changeRef{number:2}

	a.number = 1
	t.Log(a.number)
}
