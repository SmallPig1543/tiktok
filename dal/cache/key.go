package cache

import "fmt"

// UserLikesVideoKey 用户对视频的点赞id列表 采用集合
func UserLikesVideoKey(uid string) string {
	return fmt.Sprintf("like:user:%s:video", uid)
}

// VideoViewKey 视频观看量
func VideoViewKey(vid string) string {
	return fmt.Sprintf("video:view:%s", vid)
}

func VideoInfoKey(vid string) string {
	return fmt.Sprintf("video:info:%s", vid)
}

// VideoLikeKey 视频点赞数
func VideoLikeKey(vid string) string {
	return fmt.Sprintf("like:video:%s", vid)
}

// VideoCommentKey 视频评论, key是视频id， value为cid ,采用集合
func VideoCommentKey(vid string) string {
	return fmt.Sprintf("comment:video:%s", vid)
}

// CommentChildrenKey 评论的回复 key是评论id， value是子评论id 采用集合
func CommentChildrenKey(cid string) string {
	return fmt.Sprintf("comment:children:%s", cid)
}

// CommentLikeKey 统计comment的点赞次数
func CommentLikeKey(cid string) string {
	return fmt.Sprintf("comment:like:%s", cid)
}

// FollowListKey 关注列表
func FollowListKey(followerID string) string {
	return fmt.Sprintf("follower:%s", followerID)
}

// FansListKey 粉丝列表
func FansListKey(followedID string) string {
	return fmt.Sprintf("fans:%s", followedID)
}
