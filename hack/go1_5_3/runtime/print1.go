// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

// The compiler knows that a print of a value of this type
// should use printhex instead of printuint (decimal).
type hex uint64
