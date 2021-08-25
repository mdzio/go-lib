package conc

import (
	"sync"
	"time"
)

// DebouncedFunc waits the specified amount of time after a trigger call and
// then calls the specified function. Further trigger calls while waiting are
// discarded and restart the time span. A DebouncedFunc must not be copied after
// first use.
type DebouncedFunc struct {
	Dur  time.Duration
	Func func()

	tmr *time.Timer
	mtx sync.Mutex
}

func (d *DebouncedFunc) Trigger() {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	if d.tmr != nil {
		d.tmr.Reset(d.Dur)
	} else {
		d.tmr = time.AfterFunc(d.Dur, d.Func)
	}
}

func (d *DebouncedFunc) Cancel() {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	if d.tmr != nil {
		if !d.tmr.Stop() {
			<-d.tmr.C
		}
	}
}
