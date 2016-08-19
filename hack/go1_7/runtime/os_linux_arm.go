// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

const (
	_AT_PLATFORM = 15
	_AT_HWCAP    = 16

	_HWCAP_VFP   = 1 << 6
	_HWCAP_VFPv3 = 1 << 13
)
