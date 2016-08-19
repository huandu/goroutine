// Copyright 2016 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

type Hacker interface {
	Package() string        // The package to hack.
	Hack(pw *PackageWriter) // Hack a package.
}
