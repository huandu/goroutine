// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

type callbacks struct {
	lock mutex
	ctxt [cb_max]*wincallbackcontext
	n    int
}

const _LOAD_LIBRARY_SEARCH_SYSTEM32 = 0x00000800
