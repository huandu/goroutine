// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Software IEEE754 64-bit floating point.
// Only referred to (and thus linked in) by arm port
// and by tests in this directory.

package runtime

const (
	mantbits64 uint = 52
	expbits64  uint = 11
	bias64          = -1<<(expbits64-1) + 1

	nan64 uint64 = (1<<expbits64-1)<<mantbits64 + 1
	inf64 uint64 = (1<<expbits64 - 1) << mantbits64
	neg64 uint64 = 1 << (expbits64 + mantbits64)

	mantbits32 uint = 23
	expbits32  uint = 8
	bias32          = -1<<(expbits32-1) + 1

	nan32 uint32 = (1<<expbits32-1)<<mantbits32 + 1
	inf32 uint32 = (1<<expbits32 - 1) << mantbits32
	neg32 uint32 = 1 << (expbits32 + mantbits32)
)
