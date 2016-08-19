// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

const (
	stackDebug       = 0
	stackFromSystem  = 0
	stackFaultOnFree = 0
	stackPoisonCopy  = 0

	stackCache = 1
)

const (
	uintptrMask = 1<<(8*ptrSize) - 1
	poisonStack = uintptrMask & 0x6868686868686868

	stackPreempt = uintptrMask & -1314

	stackFork = uintptrMask & -1234
)

type adjustinfo struct {
	old   stack
	delta uintptr
}

// Information from the compiler about the layout of stack frames.
type bitvector struct {
	n        int32
	bytedata *uint8
}

type gobitvector struct {
	n        uintptr
	bytedata []uint8
}
