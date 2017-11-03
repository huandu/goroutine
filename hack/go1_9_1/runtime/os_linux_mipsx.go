// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux
// +build mips mipsle

package runtime

const (
	_SS_DISABLE  = 2
	_NSIG        = 128 + 1
	_SI_USER     = 0
	_SIG_BLOCK   = 1
	_SIG_UNBLOCK = 2
	_SIG_SETMASK = 3
	_RLIMIT_AS   = 6
)

type sigset [4]uint32

type rlimit struct {
	rlim_cur uintptr
	rlim_max uintptr
}
