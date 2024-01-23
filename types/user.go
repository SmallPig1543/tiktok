package types

import (
	"strconv"
	"tiktok/biz/model/user"
	"tiktok/dal/db/model"
	"tiktok/pkg/util"
)

func BuildUser(u *model.User) *user.User {
	var deletedAt, url string
	if u.DeletedAt.Valid {
		deletedAt = u.DeletedAt.Time.Format("2006-01-02 15:04:05")
	}

	if u.AvatarKey != "" {
		url, _ = util.GetURL(u.AvatarKey)
	}
	return &user.User{
		ID:        strconv.Itoa(int(u.ID)),
		Name:      u.UserName,
		AvatarURL: url,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt: deletedAt,
	}
}
