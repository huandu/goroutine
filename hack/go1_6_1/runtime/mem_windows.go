// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

const (
	_MEM_COMMIT   = 0x1000
	_MEM_RESERVE  = 0x2000
	_MEM_DECOMMIT = 0x4000
	_MEM_RELEASE  = 0x8000

	_PAGE_READWRITE = 0x0004
	_PAGE_NOACCESS  = 0x0001
)
