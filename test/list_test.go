package test

import (
	"container/list"
	"testing"
)

func TestList(t *testing.T) {
	li := list.New()
	li.PushBack(1)
	li.PushBack(2)
}

func TestClosedChan(t *testing.T) {
	c := make(chan int, 10)
	c <- 1
	c <- 2
	//close(c)
	for i := 0; i < 5; i++ {
		v, ok := <-c
		t.Logf("%v - %t", v, ok)
	}

}
