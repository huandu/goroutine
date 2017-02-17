// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Software floating point interpretation of ARM 7500 FP instructions.
// The interpretation is not bit compatible with the 7500.
// It uses true little-endian doubles, while the 7500 used mixed-endian.

package runtime

const (
	_CPSR    = 14
	_FLAGS_N = 1 << 31
	_FLAGS_Z = 1 << 30
	_FLAGS_C = 1 << 29
	_FLAGS_V = 1 << 28
)

const _FAULT = 0x80000000
