package conc

import (
	"reflect"
	"testing"
	"time"
)

func TestDaemonFunc(t *testing.T) {
	l := []int{}
	c := DaemonFunc(func(ctx Context) {
		l = append(l, 1)
		if err := ctx.Sleep(1 * time.Second); err == ErrCanceled {
			l = append(l, 2)
			return
		}
		l = append(l, 3)
	})
	c()
	if !reflect.DeepEqual(l, []int{1, 2}) {
		t.Fatal(l)
	}

	l = []int{}
	w := make(chan struct{})
	c = DaemonFunc(func(ctx Context) {
		l = append(l, 1)
		if err := ctx.Sleep(0); err == ErrCanceled {
			l = append(l, 2)
			return
		}
		l = append(l, 3)
		w <- struct{}{}
	})
	<-w
	time.Sleep(100 * time.Millisecond)
	c()
	if !reflect.DeepEqual(l, []int{1, 3}) {
		t.Fatal(l)
	}
}

func TestDaemonPool(t *testing.T) {
	p := &DaemonPool{}
	l := []int{}
	p.Run(func(ctx Context) {
		l = append(l, 1)
		p.Run(func(Context) {
			l = append(l, 2)
		})
		if err := ctx.Sleep(1 * time.Second); err == ErrCanceled {
			l = append(l, 3)
			return
		}
		l = append(l, 4)
	})
	time.Sleep(100 * time.Millisecond)
	p.Close()
	if !reflect.DeepEqual(l, []int{1, 2, 3}) {
		t.Fatal(l)
	}
}

func TestDaemonPoolNoRun(t *testing.T) {
	p := &DaemonPool{}
	p.Close()
}
