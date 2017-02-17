// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build dragonfly freebsd linux

package runtime

const (
	mutex_unlocked = 0
	mutex_locked   = 1
	mutex_sleeping = 2

	active_spin     = 4
	active_spin_cnt = 30
	passive_spin    = 1
)
