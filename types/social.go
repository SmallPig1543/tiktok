package types

import (
	"strconv"
	"tiktok/biz/model/social"
	"tiktok/dal/db/model"
	"tiktok/pkg/util"
)

func BuildSocialUser(u *model.User) *social.User {
	var url string
	if u.AvatarKey != "" {
		url, _ = util.GetURL(u.AvatarKey)
	}
	return &social.User{
		ID:        strconv.Itoa(int(u.ID)),
		Username:  u.UserName,
		AvatarURL: url,
	}
}

func BuildSocialUserList(users []*model.User) []*social.User {
	resp := make([]*social.User, 0)
	for _, data := range users {
		resp = append(resp, BuildSocialUser(data))
	}
	return resp
}
