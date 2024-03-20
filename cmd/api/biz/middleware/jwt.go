package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"tiktok/cmd/api/biz/pack"
	"tiktok/internal/errno"
	"tiktok/internal/utils"
)

func JWT() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		accessToken := string(c.GetHeader("Access-Token"))
		refreshToken := string(c.GetHeader("Refresh-Token"))
		if accessToken == "" {
			pack.RespError(c, errno.AuthorizationFailedErr)
			c.Abort()
			return
		}
		aClaims, aValid, err := utils.ParseToken(accessToken)
		bClaims, bValid, err := utils.ParseToken(refreshToken)
		//两者都过期，需要重新登录
		if !aValid && !bValid {
			pack.RespError(c, errno.AuthorizationFailedErr)
			c.Abort()
			return
		} else if bValid { //只要refresh_token没过期，就直接更新access_token
			accessToken, err = utils.GenerateAccessToken(bClaims.ID, bClaims.UserName)
			if err != nil {
				pack.RespError(c, errno.AuthorizationFailedErr)
				c.Abort()
				return
			}
		} else { //其余情况都更新
			accessToken, err = utils.GenerateAccessToken(aClaims.ID, aClaims.UserName)
			refreshToken, err = utils.GenerateRefreshToken(aClaims.ID, aClaims.UserName)
			if err != nil {
				pack.RespError(c, errno.AuthorizationFailedErr)
				c.Abort()
				return
			}
		}
		c.Header("Access-Token", accessToken)
		c.Header("Refresh-Token", refreshToken)
		claims, _, _ := utils.ParseToken(accessToken)
		ctx = pack.NewContext(ctx, &pack.UserInfo{ID: claims.ID, UserName: claims.UserName})
		c.Next(ctx)
	}
}
