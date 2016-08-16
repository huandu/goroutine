// Copyright 2015 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
)

func logFatalf(format string, args ...interface{}) {
	err := fmt.Errorf(format, args...)
	panic(err)
}

func logErrorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
}

func logTracef(format string, args ...interface{}) {
	fmt.Printf("TRACE: "+format+"\n", args...)
}
