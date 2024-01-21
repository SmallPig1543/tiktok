package service

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"sync"
	"tiktok/biz/model/user"
	"tiktok/dal/db/dao"
	"tiktok/dal/db/model"
	"tiktok/pkg/e"
	"tiktok/pkg/util"
	"tiktok/types"
)

var UserServiceOnce sync.Once

var UserServiceIns *UserService

type UserService struct {
}

func GetUserService() *UserService {
	UserServiceOnce.Do(func() {
		UserServiceIns = &UserService{}
	})
	return UserServiceIns
}

func (s *UserService) Register(ctx context.Context, req *user.RegisterRequest) (code int64, err error) {
	//先从数据库查找是否存在该用户
	userDao := dao.NewUserDao(ctx)
	_, err = userDao.FindUserByUserName(req.Username)
	//不存在就从数据库创立新数据
	if errors.Is(err, gorm.ErrRecordNotFound) {
		u := &model.User{
			Model:    gorm.Model{},
			UserName: req.Username,
		}
		//将密码加密
		if err = u.SetPassword(req.Password); err != nil {
			util.LogrusObj.Error(err)
			code = e.SetPasswordFail
			return
		}
		//将用户放入数据库
		if err = userDao.CreateUser(u); err != nil {
			util.LogrusObj.Error(err)
			code = e.ErrorDataBase
			return
		}
		return e.SUCCESS, nil
	}
	//如果存在报错返回
	return e.ErrorUserExist, errors.New("用户存在")
}

func (s *UserService) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, code int64, err error) {
	//先从数据库查找是否存在该用户
	code = e.SUCCESS
	userDao := dao.NewUserDao(ctx)
	u, err := userDao.FindUserByUserName(req.Username)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorUserNotExist
		return
	}
	// 校验密码
	if !u.CheckPassword(req.Password) {
		err = errors.New("密码错误")
		util.LogrusObj.Info(err)
		code = e.ErrorPassword
		return
	}
	//如果用户开启了otp，进行验证
	if u.TotpSecret != "" {
		ok := util.VerifyOtp(*req.Otp, u.TotpSecret)
		if !ok {
			err = errors.New("verify otp failed")
			util.LogrusObj.Info(err)
			code = e.VerifyOtpFailed
			return
		}
	}
	//生成token
	accessToken, err := util.GenerateAccessToken(u.ID, u.UserName)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.GenerateTokenFailure
		return
	}
	refreshToken, err := util.GenerateRefreshToken(u.ID, u.UserName)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.GenerateTokenFailure
		return
	}

	//构建返回结构
	User := types.BuildUser(u)
	resp = &user.LoginResponse{
		User:         User,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}
	return
}

func (s *UserService) GetUserInfo(ctx context.Context, req *user.InfoRequest) (resp *user.InfoResponse, code int64, err error) {
	//从数据库中取得数据
	userDao := dao.NewUserDao(ctx)
	res, err := userDao.FindUserByUid(req.ID)
	if err != nil {
		util.LogrusObj.Info(err)
		code = e.ErrorUserNotExist
		return
	}
	//构建返回结构
	resp = &user.InfoResponse{User: types.BuildUser(res)}
	return
}
