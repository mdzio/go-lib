package conc

import "sync"

// DaemonFunc runs the specified function as daemon. Calling the returned
// function signals that the daemon function should be canceled and waits until
// the daemon function returns.
func DaemonFunc(f func(Context)) func() {
	// size of 1 does not block go routine, if the function is completed before
	// it is cancelled.
	term := make(chan struct{}, 1)
	ctx := &context{
		done: make(chan struct{}),
	}
	go func() {
		defer func() { term <- struct{}{} }()
		f(ctx)
	}()
	return func() {
		close(ctx.done)
		<-term
	}
}

// DaemonPool runs any number of functions as daemon. All functions can be
// canceled simultaneously with Close.
type DaemonPool struct {
	once sync.Once
	ctx  *context
	wg   sync.WaitGroup
}

// Run runs the specified function as daemon. The specified function can spawn
// additional daemon functions.
func (d *DaemonPool) Run(f func(Context)) {
	d.once.Do(func() {
		d.ctx = &context{
			done: make(chan struct{}),
		}
	})
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		f(d.ctx)
	}()
}

// Close signals to all running daemon functions that they should be canceled
// and waits until all daemon functions return.
func (d *DaemonPool) Close() {
	close(d.ctx.done)
	d.wg.Wait()
}
