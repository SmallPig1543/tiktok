package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Uid         uint `json:"uid"` //发布者id
	UserName    string
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverKey    string //储存在oss中的key
	VideoKey    string
}

// ToMap 将结构体转化成map结构,放于redis
func (v *Video) ToMap() map[string]interface{} {
	var deletedAt string
	if v.DeletedAt.Valid {
		deletedAt = v.DeletedAt.Time.Format("2006-01-02 15:04:05")
	}
	return map[string]interface{}{
		"id":          v.ID,
		"uid":         v.Uid,
		"username":    v.UserName,
		"title":       v.Title,
		"description": v.Description,
		"cover_key":   v.CoverKey,
		"video_key":   v.VideoKey,
		"created_at":  v.CreatedAt.Format("2006-01-02 15:04:05"),
		"update_at":   v.UpdatedAt.Format("2006-01-02 15:04:05"),
		"delete_at":   deletedAt,
	}
}
