// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

// In addition to the 16 bits taken from the top, we can take 3 from the
// bottom, because node must be pointer-aligned, giving a total of 19 bits
// of count.
const (
	addrBits = 48
	cntBits  = 64 - addrBits + 3
)
