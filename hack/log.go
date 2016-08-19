// Copyright 2016 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
)

var flagDebug bool

func init() {
	flag.BoolVar(&flagDebug, "debug", false, "Print debug information.")
}

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

func logDebugf(format string, args ...interface{}) {
	if !flagDebug {
		return
	}

	fmt.Printf("DEBUG: "+format+"\n", args...)
}
