package types

import (
	"strconv"
	"tiktok/biz/model/user"
	"tiktok/dal/db/model"
)

func BuildUser(u *model.User) *user.User {
	deletedAt := ""
	if u.DeletedAt.Valid {
		deletedAt = u.DeletedAt.Time.Format("2006-01-02 15:04:05")
	}
	return &user.User{
		ID:        strconv.Itoa(int(u.ID)),
		Name:      u.UserName,
		AvatarURL: u.AvatarKey,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt: deletedAt,
	}
}
