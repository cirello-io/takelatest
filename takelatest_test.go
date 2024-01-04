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
