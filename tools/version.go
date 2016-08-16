// Copyright 2015 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"regexp"
	"strings"
)

var reVersion = regexp.MustCompile(`^go(\d+)(\.(\d+))+`)

type Version []string

// Parse version string like "go1.6.3".
func ParseVersion(version string) (Version, error) {
	matches := reVersion.FindStringSubmatch(version)

	if matches == nil {
		return nil, fmt.Errorf("invalid version format")
	}

	v := Version{}

	for i := 1; i < len(matches); i += 2 {
		v = append(v, matches[i])
	}

	return v, nil
}

// Return a version string like "go1.6.3".
func (v Version) String() string {
	return "go" + v.Join(".")
}

// Join version numbers with sep.
func (v Version) Join(sep string) string {
	return strings.Join(([]string)(v), sep)
}

// Compare two versions.
// Return 1 if a > b.
// Return 0 if two versions equal.
// Return -1 if a < b.
func CompareVersions(a, b Version) int {
	if a == nil && b != nil {
		return -1
	}

	if a != nil && b == nil {
		return 1
	}

	la := len(a)
	lb := len(b)

	for i := 0; i < la && i < lb; i++ {
		if c := strings.Compare(a[i], b[i]); c != 0 {
			return c
		}
	}

	if la > lb {
		return 1
	} else if la < lb {
		return -1
	} else {
		return 0
	}
}
