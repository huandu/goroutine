// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

const _PAGESIZE = 0x1000

type ureg struct {
	r0   uint32
	r1   uint32
	r2   uint32
	r3   uint32
	r4   uint32
	r5   uint32
	r6   uint32
	r7   uint32
	r8   uint32
	r9   uint32
	r10  uint32
	r11  uint32
	r12  uint32
	sp   uint32
	link uint32
	trap uint32
	psr  uint32
	pc   uint32
}

type sigctxt struct {
	u *ureg
}
