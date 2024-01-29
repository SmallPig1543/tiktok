package cache

import (
	"context"
	"strconv"
)

// AddCommentChildren 增加子评论
func AddCommentChildren(ctx context.Context, parentID, childID string) (err error) {
	err = RedisClient.SAdd(ctx, CommentChildrenKey(parentID), childID).Err()
	return
}

// CommentChildren 子评论个数
func CommentChildren(ctx context.Context, cid string) (count int64, err error) {
	count, err = RedisClient.SCard(ctx, CommentChildrenKey(cid)).Result()
	return
}

func CommentChildrenList(ctx context.Context, cid string) ([]string, error) {
	list, err := RedisClient.SMembers(ctx, CommentChildrenKey(cid)).Result()
	return list, err
}

// AddCommentLikes 对评论点赞或者取消点赞
func AddCommentLikes(ctx context.Context, cid string, actionType string) {
	if actionType == "1" {
		RedisClient.Incr(ctx, CommentLikeKey(cid))
	} else {
		RedisClient.Decr(ctx, CommentLikeKey(cid))
	}
}

func CommentLikes(ctx context.Context, cid string) int64 {
	_ = RedisClient.SetNX(ctx, CommentLikeKey(cid), 0, 0)
	countStr, _ := RedisClient.Get(ctx, CommentLikeKey(cid)).Result()
	count, _ := strconv.Atoi(countStr)
	return int64(count)
}
