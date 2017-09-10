// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"unsafe"
)

// TODO(brainman): should not need those
const (
	_NSIG = 65
)

type stdFunction unsafe.Pointer

type mOS struct {
	waitsema uintptr
}

type sigset struct{}

const (
	currentProcess = ^uintptr(0)
	currentThread  = ^uintptr(1)
)

// osRelaxMinNS indicates that sysmon shouldn't osRelax if the next
// timer is less than 60 ms from now. Since osRelaxing may reduce
// timer resolution to 15.6 ms, this keeps timer error under roughly 1
// part in 4.
const osRelaxMinNS = 60 * 1e6
