package dao

import (
	"context"
	"gorm.io/gorm"
	dao "tiktok/dal/db"
	"tiktok/dal/db/model"
)

type RelationDao struct {
	*gorm.DB
}

func NewRelationDao(ctx context.Context) *RelationDao {
	if ctx == nil {
		ctx = context.Background()
	}

	return &RelationDao{dao.NewDBClient(ctx)}
}

func (dao *RelationDao) Create(relation *model.Follow) (err error) {
	err = dao.DB.Model(&model.Follow{}).Create(relation).Error
	return
}

func (dao *RelationDao) Delete(follower, followed string) (err error) {
	err = dao.DB.Where("follower_id = ? AND followed_id = ?", follower, followed).
		Delete(&model.Follow{}).Error
	return
}
