// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Plan 9-specific system calls

package runtime

// open
const (
	_OREAD   = 0
	_OWRITE  = 1
	_ORDWR   = 2
	_OEXEC   = 3
	_OTRUNC  = 16
	_OCEXEC  = 32
	_ORCLOSE = 64
	_OEXCL   = 0x1000
)

// rfork
const (
	_RFNAMEG  = 1 << 0
	_RFENVG   = 1 << 1
	_RFFDG    = 1 << 2
	_RFNOTEG  = 1 << 3
	_RFPROC   = 1 << 4
	_RFMEM    = 1 << 5
	_RFNOWAIT = 1 << 6
	_RFCNAMEG = 1 << 10
	_RFCENVG  = 1 << 11
	_RFCFDG   = 1 << 12
	_RFREND   = 1 << 13
	_RFNOMNT  = 1 << 14
)

// notify
const (
	_NCONT = 0
	_NDFLT = 1
)

type uinptr _Plink

type tos struct {
	prof struct {
		pp    *_Plink
		next  *_Plink
		last  *_Plink
		first *_Plink
		pid   uint32
		what  uint32
	}
	cyclefreq uint64
	kcycles   int64
	pcycles   int64
	pid       uint32
	clock     uint32
}

const (
	_NSIG   = 14
	_ERRMAX = 128

	_SIGRFAULT = 2
	_SIGWFAULT = 3
	_SIGINTDIV = 4
	_SIGFLOAT  = 5
	_SIGTRAP   = 6
	_SIGPROF   = 0
	_SIGQUIT   = 0
)
