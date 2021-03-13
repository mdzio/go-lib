package conc

import (
	"reflect"
	"testing"
	"time"
)

func TestFunc(t *testing.T) {
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
}
