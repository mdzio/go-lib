package conc

import (
	"errors"
	"time"
)

var (
	// ErrCanceled is returned from functions that were canceled prematurely.
	ErrCanceled = errors.New("Canceled")
)

type Context interface {
	// Done returns a channel that is closed when work should be canceled.
	Done() <-chan struct{}

	// IsDone returns true, when work should be canceled.
	IsDone() bool

	// Sleep pauses the execution for the specified duration. If the context is
	// canceled, Sleep returns immediately ErrCanceled.
	Sleep(time.Duration) error
}

type context struct {
	done chan struct{}
}

func (c *context) Done() <-chan struct{} {
	return c.done
}

func (c *context) IsDone() bool {
	// a read on closed channel succeeds immediately
	select {
	case <-c.done:
		return true
	default:
		return false
	}
}

func (c *context) Sleep(d time.Duration) error {
	t := time.NewTimer(d)
	select {
	case <-t.C:
		// time is up
		return nil
	case <-c.done:
		// cancel received, clean up timer
		if !t.Stop() {
			<-t.C
		}
		return ErrCanceled
	}
}
