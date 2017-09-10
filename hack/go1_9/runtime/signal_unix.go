// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package runtime

const (
	_SIG_DFL uintptr = 0
	_SIG_IGN uintptr = 1
)

// gsignalStack saves the fields of the gsignal stack changed by
// setGsignalStack.
type gsignalStack struct {
	stack       stack
	stackguard0 uintptr
	stackguard1 uintptr
	stktopsp    uintptr
}
