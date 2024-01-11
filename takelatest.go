/*
Copyright 2024 U. Cirello

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package takelatest implements redux-saga effects takeLatest.
package takelatest // import "cirello.io/takelatest"

import (
	"context"
	"sync"
)

// Runner implements the takeLatest effect from redux-saga. Its zero value
// consumes the taken message and does nothing with it.
type Runner[T any] struct {
	Func func(ctx context.Context, params T)

	mu   sync.Mutex
	reqs chan T
}

// Take enqueues the execution of Func with the given parameters.
func (r *Runner[T]) Take(param T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.reqs != nil {
		r.reqs <- param
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
	r.reqs <- param
}

// Close stoppers the runner.
func (r *Runner[T]) Close() {
	r.mu.Lock()
	if r.reqs != nil {
		close(r.reqs)
		r.reqs = nil
	}
	r.mu.Unlock()
}
