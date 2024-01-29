package cache

import (
	"context"
	"errors"
)

// Follow follower指关注者， followed指被关注者
func Follow(ctx context.Context, follower, followed, actionType string) (err error) {
	if actionType == "0" { //取关
		if GetFollowers(ctx, followed) == 0 {
			return errors.New("操作失败")
		}
		if err = RedisClient.SRem(ctx, FollowListKey(follower), followed).Err(); err != nil {
			panic(err)
		}
		err = RedisClient.SRem(ctx, FansListKey(followed), follower).Err()
	} else {
		if err = RedisClient.SAdd(ctx, FollowListKey(follower), followed).Err(); err != nil {
			panic(err)
		}
		err = RedisClient.SAdd(ctx, FansListKey(followed), follower).Err()
	}
	return
}

func GetFollowers(ctx context.Context, uid string) int64 {
	count, _ := RedisClient.SCard(ctx, FollowListKey(uid)).Result()
	return count
}

func GetFollowersList(ctx context.Context, uid string) ([]string, error) {
	res, err := RedisClient.SMembers(ctx, FollowListKey(uid)).Result()
	return res, err
}

func GetFansList(ctx context.Context, uid string) ([]string, error) {
	res, err := RedisClient.SMembers(ctx, FansListKey(uid)).Result()
	return res, err
}
