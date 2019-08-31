// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package upload

import (
	"net/http"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// Upload file plugin
type (
	Upload struct {
		Path            string
		PublicURLPrefix string
		URLPrefix       string
		Save            func(c *gin.Context, files []FileInfo)
	}
	// FileInfo defined file info
	FileInfo struct {
		UID    string
		Status string
		Name   string
		Size   int64
		URL    string
	}
)

// New defined return a Upload with default property
func New() *Upload {
	up := &Upload{
		PublicURLPrefix: "/public/upload",
		URLPrefix:       "/upload",
		Path:            path.Join("assets/public/upload", ""),
	}
	return up
}

// AddOptions defined add option
func (upload *Upload) AddOptions(opts ...Option) *Upload {
	for _, v := range opts {
		v.apply(upload)
	}
	return upload
}

// Plugin for bulrush
func (upload *Upload) Plugin(router *gin.RouterGroup) {
	router.POST(upload.URLPrefix, func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		ret := make([]FileInfo, 0)
		for _, files := range form.File {
			for _, file := range files {
				filename := filepath.Base(file.Filename)
				uuid := RandString(32)
				uuidFileName := RandString(32) + string(filename[len(filename)-len(filepath.Ext(filename)):])
				size := file.Size
				if err := c.SaveUploadedFile(file, path.Join(upload.Path, uuidFileName)); err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"message": err.Error(),
					})
					return
				}
				item := FileInfo{
					UID:    uuid,
					Status: "done",
					Name:   filename,
					Size:   size,
					URL:    upload.PublicURLPrefix + "/" + uuidFileName,
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
