// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package upload

import "github.com/gin-gonic/gin"

// Option defined implement of option
type (
	Option func(*Upload) interface{}
)

// option defined implement of option
func (o Option) apply(r *Upload) *Upload { return o(r).(*Upload) }

// option defined implement of option
func (o Option) check(r *Upload) interface{} { return o(r) }

// SaveOption defined save
func SaveOption(save func(c *gin.Context, files []FileInfo)) Option {
	return Option(func(r *Upload) interface{} {
		r.save = save
		return r
	})
}
