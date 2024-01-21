package dao

import (
	"context"
	"gorm.io/gorm"
	"strconv"
	"tiktok/dal/db"
	"tiktok/dal/db/model"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &UserDao{dao.NewDBClient(ctx)}
}

func (dao *UserDao) CreateUser(user *model.User) (err error) {
	err = dao.DB.Model(&model.User{}).Create(user).Error
	return
}

func (dao *UserDao) FindUserByUid(uid string) (user *model.User, err error) {
	id, _ := strconv.Atoi(uid)
	err = dao.DB.Model(&model.User{}).Where("id=?", id).
		First(&user).Error
	return
}

func (dao *UserDao) FindUserByUserName(userName string) (user *model.User, err error) {
	err = dao.DB.Model(&model.User{}).Where("user_name=?", userName).
		First(&user).Error
	return
}

func (dao *UserDao) UpdateOtpSecret(user *model.User, secret string, url string) (err error) {
	err = dao.DB.Model(&user).Update("totp_secret", secret).Error
	if err != nil {
		return err
	}
	err = dao.DB.Model(&user).Update("totp_url", url).Error
	return err
}

func (dao *UserDao) UpdateAvatar(user *model.User, avatarURL string) (err error) {
	err = dao.DB.Model(&user).Update("avatar_file_name", avatarURL).Error
	return
}
