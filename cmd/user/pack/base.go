package pack

import (
	"tiktok/internal/errno"
	"tiktok/kitex_gen/user"
)

func BuildBaseResp(err error) *user.BaseResp {
	if err == nil {
		return &user.BaseResp{
			Code: errno.SuccessCode,
			Msg:  errno.Success.ErrMsg,
		}
	}
	Errno := errno.ConvertErr(err)
	return &user.BaseResp{
		Code: Errno.ErrCode,
		Msg:  Errno.ErrMsg,
	}
}
