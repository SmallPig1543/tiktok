package service

import (
	"context"
	"github.com/pkg/errors"
	"tiktok/cmd/user/dal/db"
)

func (s *UserService) GetSecret(ctx context.Context, uid string) (secret string, err error) {
	//从数据库中获取secret
	userDao := db.NewUserDao(ctx)
	u, err := userDao.FindUserByUid(uid)
	if err != nil {
		return "", errors.WithMessage(err, "db.FindUserByUid failed")
	}
	return u.OtpSecret, nil
}
