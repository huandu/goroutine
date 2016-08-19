// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

const (
	_GoidCacheBatch = 16
)

// freezeStopWait is a large value that freezetheworld sets
// sched.stopwait to in order to request that all Gs permanently stop.
const freezeStopWait = 0x7fffffff

type cgothreadstart struct {
	g   guintptr
	tls *uint64
	fn  unsafe.Pointer
}

// forcePreemptNS is the time slice given to a G before it is
// preempted.
const forcePreemptNS = 10 * 1000 * 1000

// To shake out latent assumptions about scheduling order,
// we introduce some randomness into scheduling decisions
// when running with the race detector.
// The need for this was made obvious by changing the
// (deterministic) scheduling order in Go 1.5 and breaking
// many poorly-written tests.
// With the randomness here, as long as the tests pass
// consistently with -race, they shouldn't have latent scheduling
// assumptions.
const randomizeScheduler = raceenabled
