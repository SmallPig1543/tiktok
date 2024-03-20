package pack

import (
	"context"
	"github.com/pkg/errors"
	"tiktok/cmd/api/biz/model/api"
	"tiktok/kitex_gen/user"
)

var userKey int

type UserInfo struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
}

func GetUserInfo(ctx context.Context) (*UserInfo, error) {
	user, ok := FromContext(ctx)
	if !ok {
		return nil, errors.New("获取用户信息错误")
	}
	return user, nil
}

func NewContext(ctx context.Context, u *UserInfo) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func FromContext(ctx context.Context) (*UserInfo, bool) {
	u, ok := ctx.Value(userKey).(*UserInfo)
	return u, ok
}

func BuildUser(u *user.User) *api.User {
	return &api.User{
		ID:        u.Id,
		Name:      u.Name,
		AvatarURL: u.Avatar,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}
}
