package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"gorm.io/gorm"
	"image/png"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"sync"
	"tiktok/biz/model/user"
	"tiktok/dal/db/dao"
	"tiktok/dal/db/model"
	"tiktok/pkg/ctl"
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
	if u.OtpSecret != "" {
		if req.Otp == nil {
			err = errors.New("缺少token")
			util.LogrusObj.Error(err)

		}
		ok := util.VerifyOtp(*req.Otp, u.OtpSecret)
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

func (s *UserService) AvatarUpload(ctx context.Context, data *multipart.FileHeader) (resp *user.AvatarUploadResponse, code int64, err error) {
	//判断文件是否是图片类型
	if ok := util.IsImage(data); !ok {
		err = errors.New("头像不符要求")
		util.LogrusObj.Error(err)
		code = e.InvalidAvatar
		return
	}
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorGetUserInfo
		return
	}
	//将图片保存到本地
	ext := strings.ToLower(path.Ext(data.Filename))
	storePath := "./static/avatar/" + u.UserName + ext
	if err = util.SaveFile(data, storePath); err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorAvatarUpload
		return
	}
	//再上传到阿里云
	key, err := util.AvatarUpload(storePath)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorOssUpload
		return
	}
	//将key保存到数据库中
	userDao := dao.NewUserDao(ctx)
	User, err := userDao.FindUserByUserName(u.UserName)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorUserNotExist
		return
	}
	if err = userDao.UpdateAvatar(User, key); err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorDataBase
		return
	}
	//构建返回结构
	resp = &user.AvatarUploadResponse{User: types.BuildUser(User)}
	return
}

func (s *UserService) GetMFAqrcode(ctx context.Context) (resp *user.GetMFAqrcodeResponse, code int64, err error) {
	//获取用户信息
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorGetUserInfo
		return
	}

	//获取otp
	key, err := util.GenerateOtp(u.UserName)
	if err != nil {
		util.LogrusObj.Info(err)
		code = e.ErrorGenerateOTP
		return
	}
	//将key保存到数据库
	userDao := dao.NewUserDao(ctx)
	User, err := userDao.FindUserByUserName(u.UserName)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorUserNotExist
		return
	}
	err = userDao.UpdateOtpSecret(User, key.Secret(), key.URL())
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorDataBase
		return
	}
	//获取qrcode
	img, _ := key.Image(200, 200)
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	storePath := "./static/qrcode/" + u.UserName + ".png"
	//将图片保存到本地
	err = os.WriteFile(storePath, buf.Bytes(), 0666)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ERROR
		return
	}
	//将图片转化为base64
	data, _ := os.ReadFile(storePath)
	baseString := base64.StdEncoding.EncodeToString(data)
	//构建返回结构
	resp = &user.GetMFAqrcodeResponse{
		Secret: key.Secret(),
		Qrcode: baseString,
	}
	return
}

func (s *UserService) MFABind(ctx context.Context, req *user.MFABindRequest) (code int64, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorGetUserInfo
		return
	}
	userDao := dao.NewUserDao(ctx)

	User, err := userDao.FindUserByUserName(u.UserName)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorDataBase
		return
	}
	//先判断secret
	if User.OtpSecret != req.Secret {
		err = errors.New("认证失败")
		util.LogrusObj.Debug(err)
		code = e.MFABindFail
		return
	}
	//再判断验证码
	if ok := util.VerifyOtp(req.Code, req.Secret); !ok {
		err = errors.New("verify otp failed")
		code = e.VerifyOtpFailed
		return
	}
	//构建返回结构
	return
}
