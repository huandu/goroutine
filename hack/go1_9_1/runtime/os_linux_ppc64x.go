// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ppc64 ppc64le

package runtime

import (
	"github.com/huandu/goroutine/hack/go1_9_1/runtime/internal/sys"
)

const (
	_PPC_FEATURE_POWER5_PLUS = 0x00020000
	_PPC_FEATURE_ARCH_2_05   = 0x00001000
	_PPC_FEATURE_POWER6_EXT  = 0x00000200
	_PPC_FEATURE_ARCH_2_06   = 0x00000100
	_PPC_FEATURE2_ARCH_2_07  = 0x80000000

	_PPC_FEATURE_HAS_ALTIVEC = 0x10000000
	_PPC_FEATURE_HAS_VSX     = 0x00000080
)

type facilities struct {
	_         [sys.CacheLineSize]byte
	isPOWER5x bool
	isPOWER6  bool
	isPOWER6x bool
	isPOWER7  bool
	isPOWER8  bool
	hasVMX    bool
	hasVSX    bool
	_         [sys.CacheLineSize]byte
}
