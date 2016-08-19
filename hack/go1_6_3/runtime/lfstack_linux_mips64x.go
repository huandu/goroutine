// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build mips64 mips64le
// +build linux

package runtime

// On mips64, Linux limits the user address space to 40 bits (see
// TASK_SIZE64 in the Linux kernel).  This has grown over time,
// so here we allow 48 bit addresses.
//
// In addition to the 16 bits taken from the top, we can take 3 from the
// bottom, because node must be pointer-aligned, giving a total of 19 bits
// of count.
const (
	addrBits = 48
	cntBits  = 64 - addrBits + 3
)
