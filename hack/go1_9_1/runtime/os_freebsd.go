// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

type mOS struct{}

// From FreeBSD's <sys/sysctl.h>
const (
	_CTL_HW      = 6
	_HW_PAGESIZE = 7
)

// Undocumented numbers from FreeBSD's lib/libc/gen/sysctlnametomib.c.
const (
	_CTL_QUERY     = 0
	_CTL_QUERY_MIB = 3
)

const (
	_CPU_SETSIZE_MAX = 32
	_CPU_CURRENT_PID = -1
)

type sigactiont struct {
	sa_handler uintptr
	sa_flags   int32
	sa_mask    sigset
}
