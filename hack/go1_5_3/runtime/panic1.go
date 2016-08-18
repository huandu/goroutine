// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

const hasLinkRegister = GOARCH == "arm" || GOARCH == "arm64" || GOARCH == "ppc64" || GOARCH == "ppc64le"
