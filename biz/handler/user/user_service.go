// Code generated by hertz generator.

package user

import (
	"context"
	"tiktok/pkg/e"
	"tiktok/service"
	"tiktok/types"

	"github.com/cloudwego/hertz/pkg/app"
	user "tiktok/biz/model/user"
)

// Register .
// @router tiktok/user/register [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.RegisterRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		types.RespError(c, e.InvalidParams)
		return
	}

	l := service.GetUserService()
	code, err := l.Register(ctx, &req)
	if err != nil {
		types.RespError(c, code)
		return
	}
	types.RespSuccess(c)
}

// Login .
// @router tiktok/user/login [POST]
func Login(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.LoginRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		types.RespError(c, e.InvalidParams)
		return
	}
	resp := new(user.LoginResponse)
	l := service.GetUserService()
	resp, code, err := l.Login(ctx, &req)
	if err != nil {
		types.RespError(c, code)
		return
	}
	types.RespSuccessWithData(c, resp)
}

// GetInfo .
// @router tiktok/user/info [GET]
func GetInfo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.InfoRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		types.RespError(c, e.InvalidParams)
		return
	}

	resp := new(user.InfoResponse)
	l := service.GetUserService()
	resp, code, err := l.GetUserInfo(ctx, &req)
	if err != nil {
		types.RespError(c, code)
		return
	}
	types.RespSuccessWithData(c, resp)
}

// AvatarUpload .
// @router tiktok/user/avatar/upload [PUT]
func AvatarUpload(ctx context.Context, c *app.RequestContext) {
	var err error
	data, err := c.FormFile("data")
	if err != nil {
		types.RespError(c, e.InvalidParams)
		return
	}
	resp := new(user.AvatarUploadResponse)
	l := service.GetUserService()
	resp, code, err := l.AvatarUpload(ctx, data)
	if err != nil {
		types.RespError(c, code)
		return
	}
	types.RespSuccessWithData(c, resp)
}

// GetMFAqrcode .
// @router tiktok/auth/mfa/qrcode [GET]
func GetMFAqrcode(ctx context.Context, c *app.RequestContext) {
	var err error

	resp := new(user.GetMFAqrcodeResponse)
	l := service.GetUserService()
	resp, code, err := l.GetMFAqrcode(ctx)
	if err != nil {
		types.RespError(c, code)
		return
	}
	types.RespSuccessWithData(c, resp)
}

// MFABind .
// @router tiktok/auth/mfa/bind [POST]
func MFABind(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.MFABindRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		types.RespError(c, e.InvalidParams)
		return
	}

	l := service.GetUserService()
	code, err := l.MFABind(ctx, &req)
	if err != nil {
		types.RespError(c, code)
		return
	}
	types.RespSuccess(c)
}
