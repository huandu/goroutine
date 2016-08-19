// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

// Prior to Android-L, logging was done through writes to /dev/log files implemented
// in kernel ring buffers. In Android-L, those /dev/log files are no longer
// accessible and logging is done through a centralized user-mode logger, logd.
//
// https://android.googlesource.com/platform/system/core/+/master/liblog/logd_write.c
type loggerType int32

const (
	unknown loggerType = iota
	legacy
	logd
)
