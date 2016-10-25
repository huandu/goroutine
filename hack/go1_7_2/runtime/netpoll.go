// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris windows

package runtime

// pollDesc contains 2 binary semaphores, rg and wg, to park reader and writer
// goroutines respectively. The semaphore can be in the following states:
// pdReady - io readiness notification is pending;
//           a goroutine consumes the notification by changing the state to nil.
// pdWait - a goroutine prepares to park on the semaphore, but not yet parked;
//          the goroutine commits to park by changing the state to G pointer,
//          or, alternatively, concurrent io notification changes the state to READY,
//          or, alternatively, concurrent timeout/close changes the state to nil.
// G pointer - the goroutine is blocked on the semaphore;
//             io notification or timeout/close changes the state to READY or nil respectively
//             and unparks the goroutine.
// nil - nothing of the above.
const (
	pdReady uintptr = 1
	pdWait  uintptr = 2
)

const pollBlockSize = 4 * 1024

// Network poller descriptor.
type pollDesc struct {
	link *pollDesc

	lock    mutex
	fd      uintptr
	closing bool
	seq     uintptr
	rg      uintptr
	rt      timer
	rd      int64
	wg      uintptr
	wt      timer
	wd      int64
	user    uint32
}

type pollCache struct {
	lock  mutex
	first *pollDesc
}
