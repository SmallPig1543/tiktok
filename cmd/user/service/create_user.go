package service

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"tiktok/cmd/user/dal/db"
	"tiktok/internal/errno"
	"tiktok/internal/utils"
	"tiktok/kitex_gen/user"
)

func (s *UserService) CreateUser(ctx context.Context, req *user.RegisterRequest) (err error) {
	userDao := db.NewUserDao(ctx)
	//判断是否已经创建过了
	_, err = userDao.FindUserByUsername(req.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		u := db.User{
			Username: req.Username,
		}
		//加密密码
		password, err := utils.SetPassword(req.Password)
		if err != nil {
			return errors.Wrap(err, "utils.SetPassword failed")
		}
		u.Password = password
		//放入数据库
		if err = userDao.CreateUser(&u); err != nil {
			return errors.Wrap(err, "db.CreateUser failed")
		}
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "db error")
	}
	return errno.UserAlreadyExistErr
}
