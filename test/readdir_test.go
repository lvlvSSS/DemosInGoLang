package test

import (
	"os"
	"testing"
)

func TestReaddirnames(t *testing.T) {
	f, _ := os.Open("/Users/nelson/Documents")
	fstate, _ := f.Stat()
	t.Logf("is dir : %t", fstate.IsDir())

	names, _ := f.Readdirnames(2)
	t.Log(len(names))
	for _, name := range names {
		t.Logf("name : %s", name)
	}
}
