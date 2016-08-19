// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

const (
	_StackSystem = goos_windows*512*ptrSize + goos_plan9*512 + goos_darwin*goarch_arm*1024

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

	_StackGuard = 640*stackGuardMultiplier + _StackSystem

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
