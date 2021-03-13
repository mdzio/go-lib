package conc

// DaemonFunc runs the specified function as daemon. Calling the returned function
// signals that the function should be canceled.
func DaemonFunc(f func(Context)) func() {
	term := make(chan struct{})
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
