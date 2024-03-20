package main

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/pkg/errors"
	"strconv"
	"tiktok/cmd/user/pack"
	"tiktok/cmd/user/service"
	"tiktok/internal/errno"
	"tiktok/internal/utils"
	user "tiktok/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	resp = user.NewRegisterResponse()
	l := service.GetUserService()
	err = l.CreateUser(ctx, req)
	if err != nil {
		klog.Error(errors.Cause(err))
		resp.BaseResp = pack.BuildBaseResp(errors.Wrap(err, "user.CreateUser failed"))
		return resp, nil
	}
	resp.BaseResp = pack.BuildBaseResp(nil)
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	resp = user.NewLoginResponse()
	l := service.GetUserService()
	u, err := l.CheckUser(ctx, req)
	if err != nil {
		klog.Error(errors.Cause(err))
		resp.BaseResp = pack.BuildBaseResp(errors.Wrap(err, "user.CheckUser failed"))
		return resp, nil
	}
	//校验通过，生成token
	accessToken, err := utils.GenerateAccessToken(strconv.Itoa(int(u.ID)), u.Username)
	if err != nil {
		klog.Error(errors.Cause(err))
		resp.BaseResp = pack.BuildBaseResp(errors.Wrap(err, "utils.GenerateAccessToken failed"))
		return resp, nil
	}
	refreshToken, err := utils.GenerateRefreshToken(strconv.Itoa(int(u.ID)), u.Username)
	if err != nil {
		klog.Error(errors.Cause(err))
		resp.BaseResp = pack.BuildBaseResp(errors.Wrap(err, "utils.GenerateRefreshToken failed"))
		return resp, nil
	}
	//pack
	resp.User = pack.BuildUser(u)
	resp.BaseResp = pack.BuildBaseResp(nil)
	resp.AccessToken = accessToken
	resp.RefreshToken = refreshToken
	return resp, nil
}

// GetInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetInfo(ctx context.Context, req *user.InfoRequest) (resp *user.InfoResponse, err error) {
	resp = user.NewInfoResponse()
	l := service.GetUserService()
	u, err := l.GetUserInfo(ctx, req.Uid)
	if err != nil {
		klog.Error(errors.Cause(err))
		resp.BaseResp = pack.BuildBaseResp(errors.Wrap(err, "utils.GetInfo failed"))
		return resp, nil
	}
	//pack
	resp.User = pack.BuildUser(u)
	return resp, nil
}

// AvatarUpload implements the UserServiceImpl interface.
func (s *UserServiceImpl) AvatarUpload(ctx context.Context, req *user.AvatarUploadRequest) (resp *user.AvatarUploadResponse, err error) {
	resp = user.NewAvatarUploadResponse()
	l := service.GetUserService()
	u, err := l.UploadAvatar(ctx, req.Data, req.Uid)
	if err != nil {
		klog.Error(err)
		resp.BaseResp = pack.BuildBaseResp(errors.Wrap(err, "user.UploadAvatar failed"))
		return resp, nil
	}
	//pack
	resp.BaseResp = pack.BuildBaseResp(nil)
	resp.User = pack.BuildUser(u)
	return
}

// GetMFAqrcode implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetMFAqrcode(ctx context.Context, req *user.GetMFAqrcodeRequest) (resp *user.GetMFAqrcodeResponse, err error) {
	resp = user.NewGetMFAqrcodeResponse()
	l := service.GetUserService()
	resp.Qrcode, resp.Secret, err = l.GetMFAqrcode(ctx, req.Uid, req.Username)
	if err != nil {
		klog.Error(errors.Cause(err))
		resp.BaseResp = pack.BuildBaseResp(errors.WithMessage(err, "user.GetMFAqrcode failed"))
		return resp, nil
	}
	return resp, nil
}

// MFABind implements the UserServiceImpl interface.
func (s *UserServiceImpl) MFABind(ctx context.Context, req *user.MFABindRequest) (resp *user.MFABindResponse, err error) {
	resp = user.NewMFABindResponse()
	//从数据库获得secret
	l := service.GetUserService()
	secret, err := l.GetSecret(ctx, req.Uid)
	if err != nil {
		klog.Error(errors.Cause(err))
		resp.BaseResp = pack.BuildBaseResp(errors.WithMessage(err, "service.GetSecret failed"))
		return resp, nil
	}
	if secret == "" {
		resp.BaseResp = pack.BuildBaseResp(errors.New("user doesn't enable mfa"))
		return resp, nil
	}
	//比较secret
	if secret != req.Secret {
		resp.BaseResp = pack.BuildBaseResp(errno.AuthorizationFailedErr)
		return resp, nil
	}
	//比较code
	if !utils.VerifyOtp(req.Code, secret) {
		resp.BaseResp = pack.BuildBaseResp(errno.AuthorizationFailedErr)
		return resp, nil
	}
	resp.BaseResp = pack.BuildBaseResp(nil)
	return resp, nil
}
