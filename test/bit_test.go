package test

import (
	"os"
	"testing"
)

func TestBit(t *testing.T) {
	var i = 34
	t.Log(i & (-i))
}

func TestGetPwd(t *testing.T) {
	t.Log(os.Getwd())
}
