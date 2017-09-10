// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"github.com/huandu/goroutine/hack/go1_9/runtime/internal/sys"
)

const (
	_HWCAP_S390_VX = 2048
)

// facilities is padded to avoid false sharing.
type facilities struct {
	_     [sys.CacheLineSize]byte
	hasVX bool
	_     [sys.CacheLineSize]byte
}
