package test

import (
	"testing"
	"time"
)

func TestMultiTimerC(t *testing.T) {
	t.Log("Starting test multiply Timer.C")

	timer := time.NewTimer(time.Second * 1)

	go func() {
		for {
			select {
			case <-timer.C:
				t.Log("goroutine 2 timer.C")
				break
			}
			//timer.Reset(time.Second)
		}

	}()
	//
	//	go func() {
	//		select {
	//		case <-timer.C:
	//			t.Log("goroutine 1 timer.C")
	//			break
	//		}
	//	}()

	t.Log("Finished")

	time.Sleep(time.Second * 5)
}

func TestHexToChar(t *testing.T) {
	var bs byte = 0xae
	first := bs & 0x0F
	second := bs >> 4
	t.Log(first)
	t.Log(second)
	t.Log(convert(first))
	t.Log(convert(second))

}

func convert(b byte) string {
	if b < 10 {
		return string(b + 48)
	}
	return string(b - 10 + 65)
}

func TestSwitchFallthrough(t *testing.T) {
	//var str := "erin"
	switch str := "nelson"; {
	case str == "nelson":
		//t.Log("nelson")
		fallthrough
	case str == "erin":
		t.Log("erin")
		//fallthrough
	default:
		t.Log("default")
	}
}
