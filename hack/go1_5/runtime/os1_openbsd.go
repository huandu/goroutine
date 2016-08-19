// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

const (
	_ESRCH       = 3
	_EAGAIN      = 35
	_EWOULDBLOCK = _EAGAIN
	_ENOTSUP     = 91

	_CLOCK_REALTIME  = 0
	_CLOCK_VIRTUAL   = 1
	_CLOCK_PROF      = 2
	_CLOCK_MONOTONIC = 3
)

const (
	sigset_none = uint32(0)
	sigset_all  = ^uint32(0)
)

// From OpenBSD's <sys/sysctl.h>
const (
	_CTL_HW  = 6
	_HW_NCPU = 3
)

type sigactiont struct {
	sa_sigaction uintptr
	sa_mask      uint32
	sa_flags     int32
}
