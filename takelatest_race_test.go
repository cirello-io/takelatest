//go:build race

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
	"testing"
	"time"
)

func TestRaceTakeClose(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var r Runner[any]
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
			r.Take(nil)
		}
	}()
	for {
		if ctx.Err() != nil {
			return
		}
		r.Close()
	}
}
