// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

const (
	_SS_DISABLE  = 4
	_SIG_BLOCK   = 1
	_SIG_UNBLOCK = 2
	_SIG_SETMASK = 3
	_NSIG        = 33
	_SI_USER     = 0

	_UC_SIGMASK = 0x01
	_UC_CPU     = 0x04

	_EAGAIN = 35
)

type mOS struct {
	waitsemacount uint32
}

const (
	_ESRCH     = 3
	_ETIMEDOUT = 60

	_CLOCK_REALTIME  = 0
	_CLOCK_VIRTUAL   = 1
	_CLOCK_PROF      = 2
	_CLOCK_MONOTONIC = 3
)

// From NetBSD's <sys/sysctl.h>
const (
	_CTL_HW      = 6
	_HW_NCPU     = 3
	_HW_PAGESIZE = 7
)

type sigactiont struct {
	sa_sigaction uintptr
	sa_mask      sigset
	sa_flags     int32
}
