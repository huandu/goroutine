// Copyright 2015 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package goroutine

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGoroutineIdConsistency(t *testing.T) {
	cnt := 10
	exit := make(chan error)

	for i := 0; i < cnt; i++ {
		go func(n int) {
			id1 := GoroutineId()
			time.Sleep(time.Duration(rand.Int63n(100)) * time.Millisecond)
			id2 := GoroutineId()

			if id1 != id2 {
				exit <- fmt.Errorf("Inconsistent goroutine id. [old:%v] [new:%v]", id1, id2)
				return
			}

			exit <- nil
		}(i)
	}

	failed := false

	for i := 0; i < cnt; i++ {
		err := <-exit

		if err != nil {
			t.Logf("Found error. [err:%v]", err)
			failed = true
		}
	}

	if failed {
		t.Fatalf("Test failed.")
	}
}
