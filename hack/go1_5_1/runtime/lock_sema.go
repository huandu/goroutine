// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin nacl netbsd openbsd plan9 solaris windows

package runtime

// This implementation depends on OS-specific implementations of
//
//	uintptr runtime·semacreate(void)
//		Create a semaphore, which will be assigned to m->waitsema.
//		The zero value is treated as absence of any semaphore,
//		so be sure to return a non-zero value.
//
//	int32 runtime·semasleep(int64 ns)
//		If ns < 0, acquire m->waitsema and return 0.
//		If ns >= 0, try to acquire m->waitsema for at most ns nanoseconds.
//		Return 0 if the semaphore was acquired, -1 if interrupted or timed out.
//
//	int32 runtime·semawakeup(M *mp)
//		Wake up mp, which is or will soon be sleeping on mp->waitsema.
//
const (
	locked uintptr = 1

	active_spin     = 4
	active_spin_cnt = 30
	passive_spin    = 1
)
