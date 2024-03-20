package pack

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"tiktok/internal/errno"
)

type Base struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

type RespWithData struct {
	Base Base `json:"base"`
	Data any  `json:"data"`
}

func RespError(c *app.RequestContext, err error) {
	Errno := errno.ConvertErr(err)
	c.JSON(consts.StatusOK, Base{
		Code: Errno.ErrCode,
		Msg:  Errno.ErrMsg,
	})
}

func RespSuccess(c *app.RequestContext) {
	Errno := errno.Success
	c.JSON(consts.StatusOK, Base{
		Code: Errno.ErrCode,
		Msg:  Errno.ErrMsg,
	})
}

func RespData(c *app.RequestContext, data any) {
	c.JSON(consts.StatusOK, RespWithData{
		Base: Base{errno.SuccessCode, "Success"},
		Data: data,
	})
}
