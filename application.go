// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package upload

import (
	"net/http"
	"path"
	"path/filepath"

	"github.com/2637309949/bulrush"
	"github.com/gin-gonic/gin"
)

// Upload file plugin
type Upload struct {
	bulrush.PNBase
	Path            string
	PublicURLPrefix string
	URLPrefix       string
	Save            func(c *gin.Context, files []map[string]interface{})
}

// New defined return a Upload with default property
func New() *Upload {
	up := &Upload{
		PublicURLPrefix: "/public/upload",
		URLPrefix:       "/upload",
		Path:            path.Join("assets/public/upload", ""),
	}
	return up
}

// Plugin for bulrush
func (upload *Upload) Plugin() bulrush.PNRet {
	return func(router *gin.RouterGroup) {
		router.POST(upload.URLPrefix, func(c *gin.Context) {
			form, err := c.MultipartForm()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				c.Abort()
			}
			ret := make([]map[string]interface{}, 0)
			for _, files := range form.File {
				for _, file := range files {
					filename := filepath.Base(file.Filename)
					uuid := bulrush.RandString(32)
					uuidFileName := bulrush.RandString(32) + string(filename[len(filename)-len(filepath.Ext(filename)):])
					if err := c.SaveUploadedFile(file, path.Join(upload.Path, uuidFileName)); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"message": err.Error(),
						})
						return
					}
					item := map[string]interface{}{
						"uid":    uuid,
						"status": "done",
						"name":   filename,
						"url":    upload.PublicURLPrefix + "/" + uuidFileName,
					}
					ret = append(ret, item)
				}
			}
			c.JSON(http.StatusOK, ret)
			if upload.Save != nil {
				upload.Save(c, ret)
			}
		})
	}
}
