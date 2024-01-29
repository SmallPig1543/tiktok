package dao

import (
	"context"
	"gorm.io/gorm"
	dao "tiktok/dal/db"
	"tiktok/dal/db/model"
)

type CommentDao struct {
	*gorm.DB
}

func NewCommentDao(ctx context.Context) *CommentDao {
	if ctx == nil {
		ctx = context.Background()
	}

	return &CommentDao{dao.NewDBClient(ctx)}
}

func (dao *CommentDao) Create(comment *model.Comment) (err error) {
	err = dao.DB.Model(&model.Comment{}).Create(comment).Error
	return
}

func (dao *CommentDao) FindCommentByID(cid uint) (comment *model.Comment, err error) {
	err = dao.DB.Model(&model.Comment{}).Where("id = ?", cid).Find(&comment).Error
	return
}

func (dao *CommentDao) FindCommentsByID(ids []string, page_num, page_size int64) (list []*model.Comment, count int64, err error) {
	err = dao.DB.Model(&model.Comment{}).Where("id in (?)", ids).Count(&count).
		Limit(int(page_size)).
		Offset(int((page_num - 1) * page_size)).
		Find(&list).Error
	return
}

func (dao *CommentDao) DeleteComment(cid, vid string) (err error) {
	err = dao.DB.Where("vid = ?", vid).Delete(&model.Comment{}, cid).Error
	return
}
