// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The standard GNU/Linux sigset type on big-endian 64-bit machines.

// +build ppc64 s390x

package runtime

const (
	_SS_DISABLE  = 2
	_NSIG        = 65
	_SI_USER     = 0
	_SIG_BLOCK   = 0
	_SIG_UNBLOCK = 1
	_SIG_SETMASK = 2
	_RLIMIT_AS   = 9
)

type sigset uint64

type rlimit struct {
	rlim_cur uintptr
	rlim_max uintptr
}
