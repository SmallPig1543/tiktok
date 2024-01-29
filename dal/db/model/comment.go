package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Vid      uint
	Uid      uint
	ParentID uint `default:"0"` //父评论id，为0表示为对视频的评论
	Content  string
}
