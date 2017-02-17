// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"github.com/huandu/goroutine/hack/go1_7_5/runtime/internal/sys"
)

const (
	_StackSystem = sys.GoosWindows*512*sys.PtrSize + sys.GoosPlan9*512 + sys.GoosDarwin*sys.GoarchArm*1024

	_StackMin = 2048

	_FixedStack0 = _StackMin + _StackSystem
	_FixedStack1 = _FixedStack0 - 1
	_FixedStack2 = _FixedStack1 | (_FixedStack1 >> 1)
	_FixedStack3 = _FixedStack2 | (_FixedStack2 >> 2)
	_FixedStack4 = _FixedStack3 | (_FixedStack3 >> 4)
	_FixedStack5 = _FixedStack4 | (_FixedStack4 >> 8)
	_FixedStack6 = _FixedStack5 | (_FixedStack5 >> 16)
	_FixedStack  = _FixedStack6 + 1

	_StackBig = 4096

	_StackGuard = 720*sys.StackGuardMultiplier + _StackSystem

	_StackSmall = 128

	_StackLimit = _StackGuard - _StackSystem - _StackSmall
)

// Goroutine preemption request.
// Stored into g->stackguard0 to cause split stack check failure.
// Must be greater than any real sp.
// 0xfffffade in hex.
const (
	_StackPreempt = uintptrMask & -1314
	_StackFork    = uintptrMask & -1234
)

const (
	stackDebug       = 0
	stackFromSystem  = 0
	stackFaultOnFree = 0
	stackPoisonCopy  = 0

	stackCache = 1
)

const (
	uintptrMask = 1<<(8*sys.PtrSize) - 1

	stackPreempt = uintptrMask & -1314

	stackFork = uintptrMask & -1234
)

type adjustinfo struct {
	old   stack
	delta uintptr
	cache pcvalueCache

	sghi uintptr
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
