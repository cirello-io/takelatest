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
