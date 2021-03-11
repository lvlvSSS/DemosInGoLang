package test

import (
	"context"
	"os"
	"testing"
)

func TestContextValue(t *testing.T) {
	ctx := context.Background()
	child1, _ := context.WithCancel(ctx)
	child2 := context.WithValue(child1, "a", "11111")
	child3, _ := context.WithCancel(child2)
	child4 := context.WithValue(child3, "b", "22222")

	t.Log(child4.Value("c"))
}

func TestOsArgs(t *testing.T) {
	t.Log(len(os.Args))
	t.Logf("%v", os.Args)
}

func TestSlice(t *testing.T) {

	s := []byte("abbbbbbb")
	t.Log(cap(s),len(s))
	s1 := append(s, 'a')
	t.Log(cap(s1),len(s1))
	s2 := append(s, 'b')
	t.Log(cap(s2),len(s2))
	//t.Log(s1, "==========", s2)
	t.Log(string(s1), "==========", string(s2))

}
