package db

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`                 // 使用uint作为ID类型，Gorm会自动处理AUTO_INCREMENT
	Username  string         `gorm:"size:100;not null"`                        // 定义varchar字段长度为100
	Password  string         `gorm:"size:100;not null"`                        // 同上
	Avatar    string         `gorm:"size:100"`                                 // 可以为空
	OtpSecret string         `gorm:"size:100;column:otp_secret"`               // 注意这里是opt_secret，看起来可能是一个笔误，如果确实是otp_secret，请做相应调整
	CreatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"` // 使用time.Time类型，并设置默认值
	UpdatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"` // 同上
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp;index"`                     // 支持软删除，使用gorm的DeletedAt类型
}

type UserDao struct {
	*gorm.DB
}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := DB
	return db.WithContext(ctx)
}

func NewUserDao(ctx context.Context) *UserDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &UserDao{NewDBClient(ctx)}
}

func (dao *UserDao) CreateUser(user *User) (err error) {
	err = dao.DB.Model(&User{}).Create(user).Error
	return
}

func (dao *UserDao) FindUserByUid(uid string) (user User, err error) {
	err = dao.DB.First(&user, uid).Error
	return user, errors.Wrapf(err, "FindUserByUid failed, id = %s, err = %v", uid, err)
}

func (dao *UserDao) FindUserByUsername(username string) (user User, err error) {
	err = dao.DB.Model(&User{}).Where("username=?", username).
		First(&user).Error
	return user, errors.Wrapf(err, "FindUserByUserName failed, username = %s, err = %v", username, err)
}

func (dao *UserDao) UpdateField(uid string, fieldName string, value interface{}) (err error) {
	updateData := map[string]interface{}{fieldName: value}
	err = dao.DB.Model(&User{}).Where("id = ?", uid).Updates(updateData).Error
	return errors.Wrapf(err, "UpdateField failed field_name = %s, value = %v", fieldName, value)
}
