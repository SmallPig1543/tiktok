package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"tiktok/cmd/user/dal/db"
	"tiktok/config"
)

func (s *UserService) UploadAvatar(ctx context.Context, data []byte, uid string) (user *db.User, err error) {
	conf := config.Config.Oss
	client, err := oss.New(conf.OssEndPoint, conf.OssAccessKeyId, conf.OssAccessKeySecret)
	// 获取存储空间
	bucket, err := client.Bucket(conf.OssBucket)
	if err != nil {
		return nil, errors.Wrap(err, "oss配置错误")
	}
	//将发送过来的文件路径转化为oss的存储路径
	objectKey := "avatar/" + uuid.Must(uuid.NewRandom()).String() + ".png"
	_ = bucket.SetObjectACL(objectKey, oss.ACLPublicReadWrite)
	err = bucket.PutObject(objectKey, bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	url := fmt.Sprintf("https://%s.%s/%s", conf.OssBucket, conf.OssEndPoint, objectKey)
	//更新数据库
	userDao := db.NewUserDao(ctx)
	if err = userDao.UpdateField(uid, "avatar", url); err != nil {
		return nil, errors.Wrap(err, "db.Update failed")
	}
	u, err := userDao.FindUserByUid(uid)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}
	return &u, nil
}
