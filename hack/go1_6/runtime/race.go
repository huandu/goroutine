// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build race

// Public race detection API, present iff build with -race.

package runtime

// private interface for the runtime
const raceenabled = true

type symbolizeContext struct {
	pc   uintptr
	fn   *byte
	file *byte
	line uintptr
	off  uintptr
	res  uintptr
}
