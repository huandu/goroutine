// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copy of math/sqrt.go, here for use by ARM softfloat.
// Modified to not use any floating point arithmetic so
// that we don't clobber any floating-point registers
// while emulating the sqrt instruction.

package runtime

const (
	float64Mask  = 0x7FF
	float64Shift = 64 - 11 - 1
	float64Bias  = 1023
	float64NaN   = 0x7FF8000000000001
	float64Inf   = 0x7FF0000000000000
	maxFloat64   = 1.797693134862315708145274237317043567981e+308
)
