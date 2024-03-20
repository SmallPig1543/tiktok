package pack

import (
	"strconv"
	"tiktok/cmd/user/dal/db"
	"tiktok/kitex_gen/user"
)

func BuildUser(u *db.User) *user.User {
	//判断是不是被删除了，如果没有返回空
	var deletedAt string
	if u.DeletedAt.Valid {
		deletedAt = u.DeletedAt.Time.Format("2006-01-02 15:04:05")
	}
	return &user.User{
		Id:        strconv.Itoa(int(u.ID)),
		Name:      u.Username,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt: deletedAt,
	}
}
