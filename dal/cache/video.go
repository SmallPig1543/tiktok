package cache

import (
	"context"
	"strconv"
)

func GetVideoInfo(ctx context.Context, vid uint) (map[string]string, bool) {
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
func AddView(ctx context.Context, vid uint) {
	RedisClient.Incr(ctx, VideoViewKey(vid))
	//排行榜
	RedisClient.ZIncrBy(ctx, "Rank", 1, strconv.Itoa(int(vid)))
}

func AddComment(ctx context.Context, vid uint) {
	RedisClient.Incr(ctx, VideoCommentKey(vid))
}

func AddLike(ctx context.Context, vid uint) {
	RedisClient.Incr(ctx, VideoLikeKey(vid))
}

// Views 获取点击量
func VideoViews(ctx context.Context, vid uint) int64 {
	_ = RedisClient.SetNX(ctx, VideoViewKey(vid), 0, 0)
	countStr, _ := RedisClient.Get(ctx, VideoViewKey(vid)).Result()
	count, _ := strconv.Atoi(countStr)
	return int64(count)
}

func VideoComments(ctx context.Context, vid uint) int64 {
	_ = RedisClient.SetNX(ctx, VideoCommentKey(vid), 0, 0)
	countStr, _ := RedisClient.Get(ctx, VideoCommentKey(vid)).Result()
	count, _ := strconv.Atoi(countStr)
	return int64(count)
}

func VideoLikes(ctx context.Context, vid uint) int64 {
	_ = RedisClient.SetNX(ctx, VideoLikeKey(vid), 0, 0)
	countStr, _ := RedisClient.Get(ctx, VideoLikeKey(vid)).Result()
	count, _ := strconv.Atoi(countStr)
	return int64(count)
}
