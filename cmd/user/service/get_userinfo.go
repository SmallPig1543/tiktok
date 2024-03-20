package service

import (
	"context"
	"github.com/pkg/errors"
	"tiktok/cmd/user/dal/db"
)

func (s *UserService) GetUserInfo(ctx context.Context, uid string) (user *db.User, err error) {
	//直接从数据库中获取
	userDao := db.NewUserDao(ctx)
	u, err := userDao.FindUserByUid(uid)
	if err != nil {
		return nil, errors.Wrap(err, "db.FindUserByUid failed")
	}
	return &u, nil
}
