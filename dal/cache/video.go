package cache

import (
	"context"
	"strconv"
)

func GetVideoInfo(ctx context.Context, vid string) (map[string]string, bool) {
	ok := RedisClient.HExists(ctx, VideoInfoKey(vid), "id")
	if !ok.Val() {
		return nil, ok.Val()
	}
	result, err := RedisClient.HGetAll(ctx, VideoInfoKey(vid)).Result()
	if err != nil {
		return nil, false
	}
	return result, true
}

func SetVideoInfo(ctx context.Context, key string, fields map[string]interface{}) (err error) {
	err = RedisClient.HMSet(ctx, key, fields).Err()
	return
}

// AddView 模拟点击
func AddView(ctx context.Context, vid string) {
	RedisClient.Incr(ctx, VideoViewKey(vid))
	//排行榜
	RedisClient.ZIncrBy(ctx, "Rank", 1, vid)
}

func AddVideoLike(ctx context.Context, vid string, actionType string) {
	if actionType == "1" {
		RedisClient.Incr(ctx, VideoLikeKey(vid))
	} else {
		RedisClient.Decr(ctx, VideoLikeKey(vid))
	}
}

// UserLikesVideo 用户对视频点赞或者取消点赞
func UserLikesVideo(ctx context.Context, uid, vid string, actionType string) (err error) {
	if actionType == "1" {
		err = RedisClient.SAdd(ctx, UserLikesVideoKey(uid), vid).Err()
	} else {
		err = RedisClient.SRem(ctx, UserLikesVideoKey(uid), vid).Err()
	}
	return
}

// UserLikesVideoList 用户点赞视频id列表
func UserLikesVideoList(ctx context.Context, uid string) (list []string, err error) {
	list, err = RedisClient.SMembers(ctx, UserLikesVideoKey(uid)).Result()
	return
}

// VideoViews 获取点击量
func VideoViews(ctx context.Context, vid string) int64 {
	_ = RedisClient.SetNX(ctx, VideoViewKey(vid), 0, 0)
	countStr, _ := RedisClient.Get(ctx, VideoViewKey(vid)).Result()
	count, _ := strconv.Atoi(countStr)
	return int64(count)
}

// AddVideoComments 增加视频评论
func AddVideoComments(ctx context.Context, vid, cid string) (err error) {
	err = RedisClient.SAdd(ctx, VideoCommentKey(vid), cid).Err()
	return
}

func VideoComments(ctx context.Context, vid string) (count int64, err error) {
	count, err = RedisClient.SCard(ctx, VideoCommentKey(vid)).Result()
	return
}

func VideoCommentsList(ctx context.Context, vid string) (list []string, err error) {
	list, err = RedisClient.SMembers(ctx, VideoCommentKey(vid)).Result()
	return
}

func VideoLikes(ctx context.Context, vid string) int64 {
	_ = RedisClient.SetNX(ctx, VideoLikeKey(vid), 0, 0)
	countStr, _ := RedisClient.Get(ctx, VideoLikeKey(vid)).Result()
	count, _ := strconv.Atoi(countStr)
	return int64(count)
}
