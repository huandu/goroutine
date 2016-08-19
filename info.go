// Copyright 2016 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package goroutine

import (
	"runtime"
	"unsafe"

	runtime1_5 "github.com/huandu/goroutine/hack/go1_5/runtime"
	runtime1_5_1 "github.com/huandu/goroutine/hack/go1_5_1/runtime"
	runtime1_5_2 "github.com/huandu/goroutine/hack/go1_5_2/runtime"
	runtime1_5_3 "github.com/huandu/goroutine/hack/go1_5_3/runtime"
	runtime1_5_4 "github.com/huandu/goroutine/hack/go1_5_4/runtime"
	runtime1_6 "github.com/huandu/goroutine/hack/go1_6/runtime"
	runtime1_6_1 "github.com/huandu/goroutine/hack/go1_6_1/runtime"
	runtime1_6_2 "github.com/huandu/goroutine/hack/go1_6_2/runtime"
	runtime1_6_3 "github.com/huandu/goroutine/hack/go1_6_3/runtime"
	runtime1_7 "github.com/huandu/goroutine/hack/go1_7/runtime"
	"github.com/huandu/goroutine/version"
)

const (
	VERSION_INVALID = iota

	VERSION1_5
	VERSION1_5_1
	VERSION1_5_2
	VERSION1_5_3
	VERSION1_5_4
	VERSION1_6
	VERSION1_6_1
	VERSION1_6_2
	VERSION1_6_3
	VERSION1_7
)

var versionCode = VERSION_INVALID

func init() {
	v := runtime.Version()
	ver, err := version.Parse(v)

	if err != nil {
		return
	}

	if ver.Compare(version.Version{"1", "5"}) == 0 {
		versionCode = VERSION1_5
	} else if ver.Compare(version.Version{"1", "5", "1"}) == 0 {
		versionCode = VERSION1_5_1
	} else if ver.Compare(version.Version{"1", "5", "2"}) == 0 {
		versionCode = VERSION1_5_2
	} else if ver.Compare(version.Version{"1", "5", "3"}) == 0 {
		versionCode = VERSION1_5_3
	} else if ver.Compare(version.Version{"1", "5", "4"}) == 0 {
		versionCode = VERSION1_5_4
	} else if ver.Compare(version.Version{"1", "6"}) == 0 {
		versionCode = VERSION1_6
	} else if ver.Compare(version.Version{"1", "6", "1"}) == 0 {
		versionCode = VERSION1_6_1
	} else if ver.Compare(version.Version{"1", "6", "2"}) == 0 {
		versionCode = VERSION1_6_2
	} else if ver.Compare(version.Version{"1", "6", "3"}) == 0 {
		versionCode = VERSION1_6_3
	} else if ver.Compare(version.Version{"1", "7"}) == 0 {
		versionCode = VERSION1_7
	}
}

func getg() unsafe.Pointer

// GoroutineId return id of current goroutine.
// It's guaranteed to be unique globally during app's life time.
func GoroutineId() int64 {
	gp := getg()

	switch versionCode {
	case VERSION1_5:
		return (*runtime1_5.Goroutine)(gp).Goid()
	case VERSION1_5_1:
		return (*runtime1_5_1.Goroutine)(gp).Goid()
	case VERSION1_5_2:
		return (*runtime1_5_2.Goroutine)(gp).Goid()
	case VERSION1_5_3:
		return (*runtime1_5_3.Goroutine)(gp).Goid()
	case VERSION1_5_4:
		return (*runtime1_5_4.Goroutine)(gp).Goid()
	case VERSION1_6:
		return (*runtime1_6.Goroutine)(gp).Goid()
	case VERSION1_6_1:
		return (*runtime1_6_1.Goroutine)(gp).Goid()
	case VERSION1_6_2:
		return (*runtime1_6_2.Goroutine)(gp).Goid()
	case VERSION1_6_3:
		return (*runtime1_6_3.Goroutine)(gp).Goid()
	case VERSION1_7:
		return (*runtime1_7.Goroutine)(gp).Goid()

	default:
		panic("unsupported go version " + runtime.Version())
	}
}
