// Copyright 2015 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package goroutine

func getg() *g

// GoroutineId return id of current goroutine.
// It's guaranteed to be unique globally during app's life time.
func GoroutineId() int64 {
    gp := getg()
    return gp.goid
}
