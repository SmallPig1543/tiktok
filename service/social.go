package service

import (
	"context"
	"strconv"
	"sync"
	"tiktok/biz/model/social"
	"tiktok/dal/cache"
	"tiktok/dal/db/dao"
	"tiktok/dal/db/model"
	"tiktok/pkg/ctl"
	"tiktok/pkg/e"
	"tiktok/pkg/util"
	"tiktok/types"
)

var SocialServiceOnce sync.Once
var SocialServiceIns *SocialService

type SocialService struct {
}

func GetSocialService() *SocialService {
	SocialServiceOnce.Do(func() {
		SocialServiceIns = &SocialService{}
	})
	return SocialServiceIns
}

func (s *SocialService) Follow(ctx context.Context, req *social.FollowRequest) (code int64, err error) {
	//获取当前用户信息
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorGetUserInfo
		return
	}
	//更新redis
	if err = cache.Follow(ctx, strconv.Itoa(int(u.ID)), req.UID, req.ActionType); err != nil {
		util.LogrusObj.Debug(err)
		code = e.ActionFails
		return
	}
	//异步更新数据库
	go func() {
		socialDao := dao.NewRelationDao(ctx)
		if req.ActionType == "0" {
			if err = socialDao.Delete(strconv.Itoa(int(u.ID)), req.UID); err != nil {
				util.LogrusObj.Fatal(err)
				panic(err)
			}
		} else {
			followedId, _ := strconv.Atoi(req.UID)
			relation := &model.Follow{
				ID:         0,
				FollowerId: u.ID,
				FollowedId: uint(followedId),
			}
			if err = socialDao.Create(relation); err != nil {
				util.LogrusObj.Fatal(err)
				panic(err)
			}
		}
	}()
	return
}

func (s *SocialService) FollowList(ctx context.Context, req *social.FollowListRequest) (resp *social.FollowListResponse, count, code int64, err error) {
	//从redis获取数据
	res, err := cache.GetFollowersList(ctx, req.UID)
	if err != nil {
		util.LogrusObj.Debug(err)
		code = e.GetListFails
		return
	}
	if len(res) == 0 { //没有结果直接返回
		return &social.FollowListResponse{}, 0, e.SUCCESS, nil
	}
	//判断一下参数
	if req.PageSize == nil {
		req.PageSize = new(int64)
		*req.PageSize = 10
	}
	if req.PageNum == nil {
		req.PageNum = new(int64)
		*req.PageNum = 0
	}
	//从数据库查找数据
	userDao := dao.NewUserDao(ctx)
	list, count, err := userDao.FindUsers(res, *req.PageNum, *req.PageSize)
	if err != nil {
		util.LogrusObj.Debug(err)
		code = e.ErrorDataBase
		return
	}
	resp = &social.FollowListResponse{Users: types.BuildSocialUserList(list)}
	return
}

func (s *SocialService) FansList(ctx context.Context, req *social.FansListRequest) (resp *social.FansListResponse, count, code int64, err error) {
	//从redis获取数据
	res, err := cache.GetFansList(ctx, req.UID)
	if err != nil {
		util.LogrusObj.Debug(err)
		code = e.GetListFails
		return
	}
	if len(res) == 0 { //没有结果直接返回
		return &social.FansListResponse{}, 0, e.SUCCESS, nil
	}
	//判断一下参数
	if req.PageSize == nil {
		req.PageSize = new(int64)
		*req.PageSize = 10
	}
	if req.PageNum == nil {
		req.PageNum = new(int64)
		*req.PageNum = 0
	}
	//从数据库查找数据
	userDao := dao.NewUserDao(ctx)
	list, count, err := userDao.FindUsers(res, *req.PageNum, *req.PageSize)
	if err != nil {
		util.LogrusObj.Debug(err)
		code = e.ErrorDataBase
		return
	}
	resp = &social.FansListResponse{Users: types.BuildSocialUserList(list)}
	return
}

func (s *SocialService) FriendsList(ctx context.Context, req *social.FriendsListRequest) (resp *social.FriendsListResponse, count, code int64, err error) {
	//获取当前用户登录状态
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorGetUserInfo
		return
	}
	//判断一下参数
	if req.PageSize == nil {
		req.PageSize = new(int64)
		*req.PageSize = 10
	}
	if req.PageNum == nil {
		req.PageNum = new(int64)
		*req.PageNum = 0
	}
	//从redis获取当前用户的关注列表和粉丝列表,进行交集运算获取好友列表
	res, err := cache.RedisClient.SInter(ctx, cache.FollowListKey(strconv.Itoa(int(u.ID))), cache.FansListKey(strconv.Itoa(int(u.ID)))).Result()
	if err != nil {
		util.LogrusObj.Debug(err)
		code = e.GetListFails
		return
	}

	if len(res) == 0 {
		return
	}
	userDao := dao.NewUserDao(ctx)
	list, count, err := userDao.FindUsers(res, *req.PageNum, *req.PageSize)
	if err != nil {
		util.LogrusObj.Debug(err)
		code = e.ErrorDataBase
		return
	}
	resp = &social.FriendsListResponse{Users: types.BuildSocialUserList(list)}
	return
}
