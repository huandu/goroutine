// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

type neverCallThisFunction struct{}

// argp used in Defer structs when there is no argp.
const _NoArgs = ^uintptr(0)
