// Package takelatest implements redux-saga effects takeLatest.
package takelatest

import (
	"context"
)

// Runner implements the takeLatest effect from redux-saga. Its zero value
// consumes the taken message and does nothing with it.
type Runner[T any] struct {
	Func func(ctx context.Context, params T)

	reqs chan T
}

// Take enqueues the execution of Func with the given parameters.
func (r *Runner[T]) Take(param T) {
	r.init()
	r.reqs <- param
}

// Close stoppers the runner.
func (r *Runner[T]) Close() {
	if r.reqs != nil {
		close(r.reqs)
		r.reqs = nil
	}
}

func (r *Runner[T]) init() {
	if r.reqs != nil {
		return
	}
	r.reqs = make(chan T)
	go func() {
		cancel := func() {}
		for params := range r.reqs {
			cancel()
			ctx, c := context.WithCancel(context.Background())
			cancel = c
			if r.Func != nil {
				go r.Func(ctx, params)
			}
		}
	}()
}
