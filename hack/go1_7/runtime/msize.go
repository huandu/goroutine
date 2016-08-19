// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Malloc small size classes.
//
// See malloc.go for overview.
//
// The size classes are chosen so that rounding an allocation
// request up to the next size class wastes at most 12.5% (1.125x).
//
// Each size class has its own page count that gets allocated
// and chopped up when new objects of the size class are needed.
// That page count is chosen so that chopping up the run of
// pages into objects of the given size wastes at most 12.5% (1.125x)
// of the memory. It is not necessary that the cutoff here be
// the same as above.
//
// The two sources of waste multiply, so the worst possible case
// for the above constraints would be that allocations of some
// size might have a 26.6% (1.266x) overhead.
// In practice, only one of the wastes comes into play for a
// given size (sizes < 512 waste mainly on the round-up,
// sizes > 512 waste mainly on the page chopping).
//
// TODO(rsc): Compute max waste for any given size.

package runtime

// divMagic holds magic constants to implement division
// by a particular constant as a shift, multiply, and shift.
// That is, given
//	m = computeMagic(d)
// then
//	n/d == ((n>>m.shift) * m.mul) >> m.shift2
//
// The magic computation picks m such that
//	d = d₁*d₂
//	d₂= 2^m.shift
//	m.mul = ⌈2^m.shift2 / d₁⌉
//
// The magic computation here is tailored for malloc block sizes
// and does not handle arbitrary d correctly. Malloc block sizes d are
// always even, so the first shift implements the factors of 2 in d
// and then the mul and second shift implement the odd factor
// that remains. Because the first shift divides n by at least 2 (actually 8)
// before the multiply gets involved, the huge corner cases that
// require additional adjustment are impossible, so the usual
// fixup is not needed.
//
// For more details see Hacker's Delight, Chapter 10, and
// http://ridiculousfish.com/blog/posts/labor-of-division-episode-i.html
// http://ridiculousfish.com/blog/posts/labor-of-division-episode-iii.html
type divMagic struct {
	shift    uint8
	mul      uint32
	shift2   uint8
	baseMask uintptr
}
