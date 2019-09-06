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
		Prefix       string
		AssetsPrefix string
		UploadPrefix string
		save         func(c *gin.Context, files []FileInfo)
	}
	// FileInfo defined file info
	FileInfo struct {
		UUID   string
		Status string
		Name   string
		Size   int64
		URL    string
	}
)

// New defined return a Upload with default property
func New() *Upload {
	up := &Upload{
		Prefix:       "/upload",
		AssetsPrefix: "/public/upload",
		UploadPrefix: path.Join("assets/public/upload", ""),
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
	router.POST(upload.Prefix, func(c *gin.Context) {
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
				fileInfo := FileInfo{
					UUID:   uuid,
					Status: "done",
					Name:   filename,
					Size:   size,
					URL:    upload.AssetsPrefix + "/" + uuidFileName,
				}
				ret = append(ret, fileInfo)
				if err := c.SaveUploadedFile(file, path.Join(upload.UploadPrefix, uuidFileName)); err != nil {
					fileInfo.Status = "error"
					if upload.save != nil {
						upload.save(c, ret)
					}
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"message": err.Error(),
					})
					// 退出处理
					return
				}
			}
		}
		c.JSON(http.StatusOK, ret)
		if upload.save != nil {
			upload.save(c, ret)
		}
	})
}
