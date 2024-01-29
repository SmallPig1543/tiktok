package model

// Follow 关注表
type Follow struct {
	ID         uint `gorm:"primarykey"`
	FollowerId uint //关注者
	FollowedId uint //被关注者
}
