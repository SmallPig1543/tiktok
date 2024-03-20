package service

import (
	"context"
	"tiktok/cmd/user/dal/db"
	"tiktok/internal/errno"
	"tiktok/internal/utils"
	"tiktok/kitex_gen/user"
)

func (s *UserService) CheckUser(ctx context.Context, req *user.LoginRequest) (user *db.User, err error) {
	userDao := db.NewUserDao(ctx)
	u, err := userDao.FindUserByUsername(req.Username)
	if err != nil {
		return nil, errno.UserNotExist
	}
	if !utils.CheckPassword(req.Password, u.Password) {
		return nil, errno.AuthorizationFailedErr
	}
	//otp检验
	if u.OtpSecret != "" {
		if req.Otp == nil {
			return nil, errno.AuthorizationFailedErr
		}
		if ok := utils.VerifyOtp(*req.Otp, u.OtpSecret); !ok {
			return nil, errno.AuthorizationFailedErr
		}
	}
	return &u, nil
}
