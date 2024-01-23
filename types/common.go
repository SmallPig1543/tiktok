package types

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"tiktok/biz/model/common"
	"tiktok/pkg/e"
)

type Response struct {
	Base common.BaseResp `json:"base"`
}

type DataList struct {
	Item  interface{} `json:"item"`
	Total int64       `json:"total"`
}

type ResponseWithData struct {
	Base common.BaseResp `json:"base"`
	Data interface{}     `json:"data"`
}

func BuildBaseResp() *common.BaseResp {
	return &common.BaseResp{
		Code: e.SUCCESS,
		Msg:  e.GetMsg(e.SUCCESS),
	}
}

func RespList(c *app.RequestContext, items interface{}, total int64) {
	resp := &ResponseWithData{
		Base: *BuildBaseResp(),
		Data: &DataList{
			Item:  items,
			Total: total,
		},
	}
	c.JSON(consts.StatusOK, resp)
}

func RespError(c *app.RequestContext, code int64) {
	c.JSON(consts.StatusOK, common.BaseResp{
		Code: code,
		Msg:  e.GetMsg(code),
	})
}

func RespSuccess(c *app.RequestContext) {
	resp := &Response{
		Base: *BuildBaseResp(),
	}
	c.JSON(consts.StatusOK, resp)
}

func RespSuccessWithData(c *app.RequestContext, data interface{}) {
	resp := &ResponseWithData{
		Base: *BuildBaseResp(),
		Data: data,
	}
	c.JSON(consts.StatusOK, resp)
}
