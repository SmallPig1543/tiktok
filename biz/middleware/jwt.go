package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"tiktok/pkg/ctl"
	"tiktok/pkg/e"
	"tiktok/pkg/util"
	"tiktok/types"
)

func JWT() app.HandlerFunc {
	f := func(ctx context.Context, c *app.RequestContext) {
		code := e.SUCCESS
		accessToken := string(c.GetHeader("Access-Token"))
		refreshToken := string(c.GetHeader("Refresh-Token"))
		if accessToken == "" {
			code = http.StatusNotFound
			types.RespError(c, code)
			c.Abort()
			return
		}
		aClaims, aValid, err := util.ParseToken(accessToken)
		if err != nil {
			code = e.ParseTokenFailure
			util.LogrusObj.Error(err)
			types.RespError(c, code)
			c.Abort()
			return
		}
		bClaims, bValid, err := util.ParseToken(refreshToken)
		if err != nil {
			code = e.ParseTokenFailure
			util.LogrusObj.Error(err)
			types.RespError(c, code)
			c.Abort()
			return
		}
		//两者都过期，需要重新登录
		if !aValid && !bValid {
			code = e.ErrorTokenTimeout
			types.RespError(c, code)

			c.Abort()
			return
		} else if bValid { //只要refresh_token没过期，就直接更新access_token
			accessToken, err = util.GenerateAccessToken(aClaims.ID, aClaims.UserName)
			if err != nil {
				code = e.TokenGeneratedFail
				util.LogrusObj.Error(err)
				types.RespError(c, code)
				c.Abort()
				return
			}
		} else { //其余情况都更新
			accessToken, err = util.GenerateAccessToken(aClaims.ID, aClaims.UserName)
			refreshToken, err = util.GenerateRefreshToken(bClaims.ID, bClaims.UserName)
			if err != nil {
				code = e.TokenGeneratedFail
				util.LogrusObj.Error(err)
				types.RespError(c, code)
				c.Abort()
				return
			}
		}
		c.Header("Access-Token", accessToken)
		c.Header("Refresh-Token", refreshToken)
		ctx = ctl.NewContext(ctx, &ctl.UserInfo{ID: aClaims.ID, UserName: aClaims.UserName})
		c.Next(ctx)
	}
	return f
}
