// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package upload

import "math/rand"

// RandString defined gen random string
func RandString(n int) string {
	const seeds = "abcdefghijklmnopqrstuvwxyz1234567890"
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = seeds[rand.Intn(len(seeds))]
	}
	return string(bytes)
}
