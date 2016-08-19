// Copyright 2016 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package runtime

// Goroutine is the internal type represents a goroutine.
type Goroutine g

// Get goid.
func (g *Goroutine) Goid() int64 {
	return g.goid
}
