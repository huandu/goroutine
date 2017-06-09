// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"unsafe"
)

// A gcSweepBuf is a set of *mspans.
//
// gcSweepBuf is safe for concurrent push operations *or* concurrent
// pop operations, but not both simultaneously.
type gcSweepBuf struct {
	spineLock mutex
	spine     unsafe.Pointer
	spineLen  uintptr
	spineCap  uintptr

	index uint32
}

const (
	gcSweepBlockEntries    = 512
	gcSweepBufInitSpineCap = 256
)

type gcSweepBlock struct {
	spans [gcSweepBlockEntries]*mspan
}
