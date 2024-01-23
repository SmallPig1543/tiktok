package util

import (
	"github.com/h2non/filetype"
	"io"
	"mime/multipart"
	"os"
)

func IsImage(data *multipart.FileHeader) bool {
	file, _ := data.Open()
	buffer := make([]byte, 261)
	_, err := file.Read(buffer)
	if err != nil {
		LogrusObj.Debug(err)
		return false
	}
	return filetype.IsImage(buffer)
}

func IsVideo(data *multipart.FileHeader) bool {
	file, _ := data.Open()
	buffer := make([]byte, 261) // 读取足够多的字节以便确定文件类型

	_, err := file.Read(buffer)
	if err != nil {
		LogrusObj.Debug(err)
		return false
	}
	return filetype.IsVideo(buffer)
}

func SaveFile(data *multipart.FileHeader, storePath string) (err error) {
	//打开本地文件
	dist, err := os.OpenFile(storePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer func(dist *os.File) {
		_ = dist.Close()
	}(dist)

	src, err := data.Open()
	if err != nil {
		return
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)
	_, err = io.Copy(dist, src)
	return
}
