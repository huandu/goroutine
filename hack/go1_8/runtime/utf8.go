// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

// Numbers fundamental to the encoding.
const (
	runeError = '\uFFFD'
	runeSelf  = 0x80
	maxRune   = '\U0010FFFF'
)

// Code points in the surrogate range are not valid for UTF-8.
const (
	surrogateMin = 0xD800
	surrogateMax = 0xDFFF
)

const (
	t1 = 0x00
	tx = 0x80
	t2 = 0xC0
	t3 = 0xE0
	t4 = 0xF0
	t5 = 0xF8

	maskx = 0x3F
	mask2 = 0x1F
	mask3 = 0x0F
	mask4 = 0x07

	rune1Max = 1<<7 - 1
	rune2Max = 1<<11 - 1
	rune3Max = 1<<16 - 1

	locb = 0x80
	hicb = 0xBF
)
