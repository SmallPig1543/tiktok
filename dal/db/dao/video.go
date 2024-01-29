package dao

import (
	"context"
	"gorm.io/gorm"
	dao "tiktok/dal/db"
	"tiktok/dal/db/model"
	"time"
)

type VideoDao struct {
	*gorm.DB
}

func NewVideoDao(ctx context.Context) *VideoDao {
	if ctx == nil {
		ctx = context.Background()
	}

	return &VideoDao{dao.NewDBClient(ctx)}
}

func (dao *VideoDao) CreateVideo(video *model.Video) (err error) {
	err = dao.DB.Model(&model.Video{}).Create(video).Error
	return
}

func (dao *VideoDao) VideoList(uid uint, page_num, page_size int) (list []*model.Video, count int64, err error) {
	offset := (page_num - 1) * page_size //偏移量
	err = dao.DB.Model(&model.Video{}).Where("uid=?", uid).Count(&count).
		Limit(page_size).
		Offset(offset).
		Find(&list).Error
	return
}

func (dao *VideoDao) VideoFeed(time time.Time) (list []*model.Video, count int64, err error) {
	err = dao.DB.Model(&model.Video{}).Where("created_at > ?", time).Count(&count).Find(&list).Error
	return
}

func (dao *VideoDao) FindVideoByVid(vid string) (video *model.Video, err error) {
	err = dao.DB.Model(&model.Video{}).Where("id = ?", vid).Find(&video).Error
	return
}

func (dao *VideoDao) FindVideosById(ids []string, page_num, page_size int64) (videos []*model.Video, count int64, err error) {
	err = dao.DB.Model(&model.Video{}).Where("id in (?)", ids).
		Count(&count).
		Limit(int(page_size)).
		Offset(int((page_num - 1) * page_size)).
		Find(&videos).Error
	return
}
