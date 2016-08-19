// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

type mOS struct {
	waitsema      int32
	waitsemacount int32
	waitsemalock  int32
}

type sigset struct{}
