// Copyright 2015 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package goroutine

import (
	//"unsafe"
	"fmt"
)

// GoroutineId return id of current goroutine.
// It's guaranteed to be unique globally during app's life time.
func GoroutineId() int64 {
	gp := getg()
	return gp.goid
}

func hack_goexit() {
	fmt.Println("Goroutine ID:", GoroutineId())
    real_goexit(_PCQuantum)
}
