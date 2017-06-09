// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

type mOS struct{}

const (
	_FUTEX_WAIT = 0
	_FUTEX_WAKE = 1
)

// Clone, the Linux rfork.
const (
	_CLONE_VM             = 0x100
	_CLONE_FS             = 0x200
	_CLONE_FILES          = 0x400
	_CLONE_SIGHAND        = 0x800
	_CLONE_PTRACE         = 0x2000
	_CLONE_VFORK          = 0x4000
	_CLONE_PARENT         = 0x8000
	_CLONE_THREAD         = 0x10000
	_CLONE_NEWNS          = 0x20000
	_CLONE_SYSVSEM        = 0x40000
	_CLONE_SETTLS         = 0x80000
	_CLONE_PARENT_SETTID  = 0x100000
	_CLONE_CHILD_CLEARTID = 0x200000
	_CLONE_UNTRACED       = 0x800000
	_CLONE_CHILD_SETTID   = 0x1000000
	_CLONE_STOPPED        = 0x2000000
	_CLONE_NEWUTS         = 0x4000000
	_CLONE_NEWIPC         = 0x8000000

	cloneFlags = _CLONE_VM |
		_CLONE_FS |
		_CLONE_FILES |
		_CLONE_SIGHAND |
		_CLONE_THREAD
)

const (
	_AT_NULL   = 0
	_AT_PAGESZ = 6
	_AT_HWCAP  = 16
	_AT_RANDOM = 25
	_AT_HWCAP2 = 26
)
