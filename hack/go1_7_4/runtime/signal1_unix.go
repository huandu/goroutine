// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package runtime

const (
	_SIG_DFL uintptr = 0
	_SIG_IGN uintptr = 1
)

// sigmask represents a general signal mask compatible with the GOOS
// specific sigset types: the signal numbered x is represented by bit x-1
// to match the representation expected by sigprocmask.
type sigmask [(_NSIG + 31) / 32]uint32
