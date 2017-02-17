/*
 * The authors of this software are Rob Pike and Ken Thompson.
 *              Copyright (c) 2002 by Lucent Technologies.
 *              Portions Copyright 2009 The Go Authors. All rights reserved.
 * Permission to use, copy, modify, and distribute this software for any
 * purpose without fee is hereby granted, provided that this entire notice
 * is included in all copies of any software which is or includes a copy
 * or modification of this software and in all copies of the supporting
 * documentation for such software.
 * THIS SOFTWARE IS BEING PROVIDED "AS IS", WITHOUT ANY EXPRESS OR IMPLIED
 * WARRANTY.  IN PARTICULAR, NEITHER THE AUTHORS NOR LUCENT TECHNOLOGIES MAKE ANY
 * REPRESENTATION OR WARRANTY OF ANY KIND CONCERNING THE MERCHANTABILITY
 * OF THIS SOFTWARE OR ITS FITNESS FOR ANY PARTICULAR PURPOSE.
 */

/*
 * This code is copied, with slight editing due to type differences,
 * from a subset of ../lib9/utf/rune.c [which no longer exists]
 */

package runtime

const (
	bit1 = 7
	bitx = 6
	bit2 = 5
	bit3 = 4
	bit4 = 3
	bit5 = 2

	t1 = ((1 << (bit1 + 1)) - 1) ^ 0xFF
	tx = ((1 << (bitx + 1)) - 1) ^ 0xFF
	t2 = ((1 << (bit2 + 1)) - 1) ^ 0xFF
	t3 = ((1 << (bit3 + 1)) - 1) ^ 0xFF
	t4 = ((1 << (bit4 + 1)) - 1) ^ 0xFF
	t5 = ((1 << (bit5 + 1)) - 1) ^ 0xFF

	rune1 = (1 << (bit1 + 0*bitx)) - 1
	rune2 = (1 << (bit2 + 1*bitx)) - 1
	rune3 = (1 << (bit3 + 2*bitx)) - 1
	rune4 = (1 << (bit4 + 3*bitx)) - 1

	maskx = (1 << bitx) - 1
	testx = maskx ^ 0xFF

	runeerror = 0xFFFD
	runeself  = 0x80

	surrogateMin = 0xD800
	surrogateMax = 0xDFFF

	bad = runeerror

	runemax = 0x10FFFF
)
