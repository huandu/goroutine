// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

const (
	_NSIG        = 33
	_SI_USER     = 0
	_SS_DISABLE  = 4
	_RLIMIT_AS   = 10
	_SIG_BLOCK   = 1
	_SIG_UNBLOCK = 2
	_SIG_SETMASK = 3
)

type mOS struct{}

const stackSystem = 0

// From DragonFly's <sys/sysctl.h>
const (
	_CTL_HW  = 6
	_HW_NCPU = 3
)

type sigactiont struct {
	sa_sigaction uintptr
	sa_flags     int32
	sa_mask      sigset
}
