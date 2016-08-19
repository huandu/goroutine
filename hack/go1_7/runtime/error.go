// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

// The Error interface identifies a run time error.
type Error interface {
	error

	RuntimeError()
}

// A TypeAssertionError explains a failed type assertion.
type TypeAssertionError struct {
	interfaceString string
	concreteString  string
	assertedString  string
	missingMethod   string
}

// An errorString represents a runtime error described by a single string.
type errorString string

// plainError represents a runtime error described a string without
// the prefix "runtime error: " after invoking errorString.Error().
// See Issue #14965.
type plainError string

type stringer interface {
	String() string
}
