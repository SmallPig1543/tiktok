package util

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"tiktok/config"
	"time"
)

var conf config.Oss

func OssInit() {
	c := config.Config.Oss
	conf = c
}
func GetURL(key string) (string, error) {
	client, err := oss.New(conf.OssEndPoint, conf.OssAccessKeyId, conf.OssAccessKeySecret)
	if err != nil {
		return "", errors.New("oss配置错误")
	}
	bucket, err := client.Bucket(conf.OssBucket)
	if err != nil {
		return "", errors.New("oss配置错误")
	}
	signedGetURL, _ := bucket.SignURL(key, oss.HTTPGet, 300)
	return signedGetURL, nil
}

func CoverUpload(FileName string) (string, error) {
	client, err := oss.New(conf.OssEndPoint, conf.OssAccessKeyId, conf.OssAccessKeySecret)
	if err != nil {
		return "", errors.New("oss配置错误")
	}
	// 获取存储空间
	bucket, err := client.Bucket(conf.OssBucket)
	if err != nil {
		return "", errors.New("oss配置错误")
	}
	// 获取扩展名
	ext := filepath.Ext(FileName)
	//将发送过来的文件路径转化为oss的存储路径
	objectKey := "cover/" + uuid.Must(uuid.NewRandom()).String() + ext
	err = bucket.PutObjectFromFile(objectKey, FileName)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("文件有误")
	}
	return objectKey, nil
}

func AvatarUpload(FileName string) (string, error) {
	client, err := oss.New(conf.OssEndPoint, conf.OssAccessKeyId, conf.OssAccessKeySecret)
	if err != nil {
		return "", errors.New("oss配置错误")
	}
	// 获取存储空间
	bucket, err := client.Bucket(conf.OssBucket)
	if err != nil {
		return "", errors.New("oss配置错误")
	}
	// 获取扩展名
	ext := filepath.Ext(FileName)
	//将发送过来的文件路径转化为oss的存储路径
	objectKey := "avatar/" + uuid.Must(uuid.NewRandom()).String() + ext
	err = bucket.PutObjectFromFile(objectKey, FileName)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("文件有误")
	}
	return objectKey, nil
}

func VideoUpload(localFileName string) (string, error) {
	client, err := oss.New(conf.OssEndPoint, conf.OssAccessKeyId, conf.OssAccessKeySecret)
	if err != nil {
		return "", errors.New("oss配置错误")
	}
	// 获取存储空间
	bucket, err := client.Bucket(conf.OssBucket)
	if err != nil {
		return "", errors.New("oss配置错误")
	}
	// 获取扩展名
	ext := filepath.Ext(localFileName)
	//将发送过来的文件路径转化为oss的存储路径
	objectName := "video/" + uuid.Must(uuid.NewRandom()).String() + ext
	// 将本地文件分片，且分片数量指定为3。
	chunks, err := oss.SplitFileByPartNum(localFileName, 3)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	fd, err := os.Open(localFileName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	defer fd.Close()
	// 指定过期时间。
	expires := time.Date(2049, time.January, 10, 23, 0, 0, 0, time.UTC)
	options := []oss.Option{
		oss.MetadataDirective(oss.MetaReplace),
		oss.Expires(expires),
	}
	//初始化一个分片上传事件。
	imur, err := bucket.InitiateMultipartUpload(objectName, options...)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	// 上传分片。
	var parts []oss.UploadPart
	for _, chunk := range chunks {
		fd.Seek(chunk.Offset, os.SEEK_SET)
		// 调用UploadPart方法上传每个分片。
		part, err := bucket.UploadPart(imur, fd, chunk.Size, chunk.Number)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}
		parts = append(parts, part)
	}
	objectAcl := oss.ObjectACL(oss.ACLDefault)
	//完成分片上传。
	cmur, err := bucket.CompleteMultipartUpload(imur, parts, objectAcl)
	if err != nil {
		LogrusObj.Error("上传失败")
		os.Exit(-1)
	}
	return cmur.Key, nil
}
