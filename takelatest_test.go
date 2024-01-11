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

package takelatest

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Example_debouncedTimeout() {
	// The latest Take() call will be executed after previous calls are canceled due to the timeout.
	done := make(chan struct{})
	r := &Runner[int]{
		Func: func(ctx context.Context, i int) {
			time.Sleep(1 * time.Second)
			if ctx.Err() != nil {
				return
			}
			fmt.Print(i)
			close(done)
		},
	}
	defer r.Close()
	r.Take(1)
	r.Take(2)
	r.Take(3)
	r.Take(4)
	r.Take(5)
	<-done
	// Output:
	// 5
}

func TestTSRBug(t *testing.T) {
	// The runner must not Terminate-and-Stay-Running.
	var observed int
	done := make(chan struct{})
	time.AfterFunc(1*time.Second, func() {
		close(done)
	})
	r := Runner[any]{
		Func: func(ctx context.Context, _ any) {
			select {
			case <-done:
				observed = 1
			case <-ctx.Done():
			}
		},
	}
	r.Take(nil)
	r.Close()
	<-done
	if observed == 1 {
		t.Fatal("trailing execution did not stop")
	}
}

func TestChangingFunc(t *testing.T) {
	// The runner must not allow Func to be updated on the fly
	observed := make(chan struct{})
	r := Runner[any]{}
	r.Take(nil)
	r.Func = func(context.Context, any) {
		observed <- struct{}{}
	}
	r.Take(nil)
	select {
	case <-observed:
		t.Fatal("runner must not be reconfigurable once started")
	case <-time.After(1 * time.Second):
	}
	r.Close()
	r.Take(nil)
	select {
	case <-observed:
	case <-time.After(1 * time.Second):
		t.Fatal("runner must be reconfigurable once closed")
	}
}
