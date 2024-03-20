package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/pkg/errors"
	"image/png"
	"os"
	"tiktok/cmd/user/dal/db"
	"tiktok/internal/utils"
)

func (s *UserService) GetMFAqrcode(ctx context.Context, uid, username string) (baseString string, secret string, err error) {
	//获取otp
	key, err := utils.GenerateOtp(username)
	if err != nil {
		return "", "", errors.Wrap(err, "GenerateOtp failed")
	}
	//将key数据保存到数据库
	userDao := db.NewUserDao(ctx)
	if err = userDao.UpdateField(uid, "otp_secret", key.Secret()); err != nil {
		return "", "", errors.Wrap(err, "db.UpdateField error")
	}

	//获取qrcode
	img, _ := key.Image(200, 200)
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	//将图片保存到本地
	storePath := "/home/smallpig/projects/tiktok-v2/cmd/user/static/qrcode/" + username + ".png"
	err = os.WriteFile(storePath, buf.Bytes(), 0666)
	if err != nil {
		return "", "", errors.Wrap(err, "open file failed")
	}
	//将图片转化为base64编码
	data, _ := os.ReadFile(storePath)
	baseString = base64.StdEncoding.EncodeToString(data)
	return baseString, key.Secret(), nil
}
