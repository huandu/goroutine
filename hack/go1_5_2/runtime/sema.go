// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Semaphore implementation exposed to Go.
// Intended use is provide a sleep and wakeup
// primitive that can be used in the contended case
// of other synchronization primitives.
// Thus it targets the same goal as Linux's futex,
// but it has much simpler semantics.
//
// That is, don't think of these as semaphores.
// Think of them as a way to implement sleep and wakeup
// such that every sleep is paired with a single wakeup,
// even if, due to races, the wakeup happens before the sleep.
//
// See Mullender and Cox, ``Semaphores in Plan 9,''
// http://swtch.com/semaphore.pdf

package runtime

type semaRoot struct {
	lock  mutex
	head  *sudog
	tail  *sudog
	nwait uint32
}

// Prime to not correlate with any user patterns.
const semTabSize = 251

// Synchronous semaphore for sync.Cond.
type syncSema struct {
	lock mutex
	head *sudog
	tail *sudog
}
