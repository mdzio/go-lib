package conc

// DaemonFunc runs the specified function as daemon. Calling the returned function
// signals that the function should be canceled.
func DaemonFunc(f func(Context)) func() {
	// size of 1 does not block go routine, if the function is completed before
	// it is cancelled.
	term := make(chan struct{}, 1)
	ctx := &context{
		done: make(chan struct{}),
	}
	go func() {
		f(ctx)
		term <- struct{}{}
	}()
	return func() {
		close(ctx.done)
		<-term
	}
}
