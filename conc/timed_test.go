package conc

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestDebounceFunc(t *testing.T) {
	var cnt int32
	df := DebouncedFunc{
		Dur: 50 * time.Millisecond,
		Func: func() {
			atomic.AddInt32(&cnt, 1)
		},
	}

	if atomic.LoadInt32(&cnt) != 0 {
		t.Fatal()
	}
	df.Trigger()
	if atomic.LoadInt32(&cnt) != 0 {
		t.Fatal()
	}
	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&cnt) != 1 {
		t.Fatal()
	}
	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&cnt) != 1 {
		t.Fatal()
	}
	df.Trigger()
	df.Cancel()
	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&cnt) != 1 {
		t.Fatal()
	}
}
