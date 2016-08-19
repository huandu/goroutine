// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

const (
	currentProcess = ^uintptr(0)
	currentThread  = ^uintptr(1)
)

// Described in http://www.dcl.hpi.uni-potsdam.de/research/WRK/2007/08/getting-os-information-the-kuser_shared_data-structure/
type _KSYSTEM_TIME struct {
	LowPart   uint32
	High1Time int32
	High2Time int32
}

const (
	_INTERRUPT_TIME = 0x7ffe0008
	_SYSTEM_TIME    = 0x7ffe0014
)
