package rpc

import (
	"context"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/pkg/errors"
	"tiktok/config"
	"tiktok/internal/errno"
	"tiktok/kitex_gen/user"
	"tiktok/kitex_gen/user/userservice"
)

func InitUserRPC() {
	conf := config.Config
	r, err := etcd.NewEtcdResolver([]string{conf.EtcdHost + ":" + conf.EtcdPort})
	if err != nil {
		panic(err)
	}
	userClient, err = userservice.NewClient("user", client.WithResolver(r))
	if err != nil {
		panic(err)
	}
}

func Register(ctx context.Context, req *user.RegisterRequest) (err error) {
	resp, err := userClient.Register(ctx, req)
	if err != nil {
		return errors.Wrap(err, "api.rpc.user Register failed")
	}
	if resp.BaseResp.Code != errno.SuccessCode {
		return errno.NewErrNo(resp.BaseResp.Code, resp.BaseResp.Msg)
	}
	return
}

func Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	resp, err = userClient.Login(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "api.rpc.user Login failed")
	}
	if resp.BaseResp.Code != errno.SuccessCode {
		return nil, errno.NewErrNo(resp.BaseResp.Code, resp.BaseResp.Msg)
	}
	return
}

func GetUserInfo(ctx context.Context, req *user.InfoRequest) (user *user.User, err error) {
	resp, err := userClient.GetInfo(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "api.rpc.user GetUserInfo failed")
	}
	if resp.BaseResp.Code != errno.SuccessCode {
		return nil, errno.NewErrNo(resp.BaseResp.Code, resp.BaseResp.Msg)
	}
	return resp.User, nil
}

func UploadAvatar(ctx context.Context, req *user.AvatarUploadRequest) (user *user.User, err error) {
	resp, err := userClient.AvatarUpload(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "api.rpc.user UploadAvatar failed")
	}
	if resp.BaseResp.Code != errno.SuccessCode {
		return nil, errno.NewErrNo(resp.BaseResp.Code, resp.BaseResp.Msg)
	}
	return resp.User, nil
}

func GetMFAqrcode(ctx context.Context, req *user.GetMFAqrcodeRequest) (qrcode, secret string, err error) {
	resp, err := userClient.GetMFAqrcode(ctx, req)
	if err != nil {
		return "", "", errors.Wrap(err, "api.rpc.user GetMFAqrcode failed")
	}
	if resp.BaseResp.Code != errno.SuccessCode {
		return "", "", errno.NewErrNo(resp.BaseResp.Code, resp.BaseResp.Msg)
	}
	return resp.Qrcode, resp.Secret, nil
}

func MFABind(ctx context.Context, req *user.MFABindRequest) (err error) {
	resp, err := userClient.MFABind(ctx, req)
	if err != nil {
		return errors.Wrap(err, "api.rpc.user MFABind failed")
	}
	if resp.BaseResp.Code != errno.SuccessCode {
		return errno.NewErrNo(resp.BaseResp.Code, resp.BaseResp.Msg)
	}
	return nil
}
