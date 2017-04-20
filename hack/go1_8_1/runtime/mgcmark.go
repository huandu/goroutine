// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Garbage collector: marking and scanning

package runtime

const (
	fixedRootFinalizers = iota
	fixedRootFreeGStacks
	fixedRootCount

	rootBlockBytes = 256 << 10

	rootBlockSpans = 8 * 1024

	maxObletBytes = 128 << 10

	idleCheckThreshold = 100000
)

type gcDrainFlags int

const (
	gcDrainUntilPreempt gcDrainFlags = 1 << iota
	gcDrainNoBlock
	gcDrainFlushBgCredit
	gcDrainIdle

	gcDrainBlock gcDrainFlags = 0
)
