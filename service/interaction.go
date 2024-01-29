package service

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"tiktok/biz/model/interaction"
	"tiktok/dal/cache"
	"tiktok/dal/db/dao"
	"tiktok/dal/db/model"
	"tiktok/pkg/ctl"
	"tiktok/pkg/e"
	"tiktok/pkg/util"
	"tiktok/types"
)

var InteractionServiceOnce sync.Once

var InteractionServiceIns *InteractionService

type InteractionService struct {
}

func GetInteractionService() *InteractionService {
	InteractionServiceOnce.Do(func() {
		InteractionServiceIns = &InteractionService{}
	})
	return InteractionServiceIns
}

func (s *InteractionService) Like(ctx context.Context, req *interaction.LikeRequest) (code int64, err error) {
	//先判断点赞类型
	// 如果两者都没有，直接返回
	if req.Vid == nil && req.Cid == nil {
		err = errors.New("缺少参数")
		util.LogrusObj.Info(err)
		code = e.ParamsMissing
		return
	}
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorGetUserInfo
		return
	}
	//如果是点赞视频
	if req.Vid != nil {
		//向redis更新数据
		//用户对视频的点赞列表
		if err = cache.UserLikesVideo(ctx, strconv.Itoa(int(u.ID)), *req.Vid, req.ActionType); err != nil {
			util.LogrusObj.Debug(err)
			code = e.ErrorRedis
			return
		}
		//视频的点赞量
		cache.AddVideoLike(ctx, *req.Vid, req.ActionType)
		//视频播放量默认增加
		cache.AddView(ctx, *req.Vid)
	} else {
		//点赞评论
		//增加点赞量
		cache.AddCommentLikes(ctx, *req.Cid, req.ActionType)
	}
	return e.SUCCESS, nil
}

func (s *InteractionService) LikeList(ctx context.Context, req *interaction.LikeListRequest) (resp *interaction.LikeListResponse, count, code int64, err error) {
	//从redis获取点赞列表
	list, err := cache.UserLikesVideoList(ctx, req.UID)
	if err != nil {
		util.LogrusObj.Debug(err)
		code = e.ErrorRedis
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
	//从数据库中取得数据
	videoDao := dao.NewVideoDao(ctx)

	res, count, err := videoDao.FindVideosById(list, *req.PageNum, *req.PageSize)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorDataBase
		return
	}
	//构建返回结构
	resp = &interaction.LikeListResponse{Videos: types.BuildVideoList(ctx, res)}
	return
}

func (s *InteractionService) Comment(ctx context.Context, req *interaction.CommentPublishRequest) (resp *interaction.CommentPublishResponse, code int64, err error) {
	//判断参数
	if req.Vid == nil && req.Cid == nil {
		err = errors.New("缺少参数")
		util.LogrusObj.Info(err)
		code = e.ParamsMissing
		return
	}
	//获取用户信息
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorGetUserInfo
		return
	}
	commentDao := dao.NewCommentDao(ctx)
	var comment *model.Comment
	//评论视频
	if req.Vid != nil && req.Cid == nil {
		//向数据库存储数据
		vid, _ := strconv.ParseUint(*req.Vid, 10, 0)
		comment = &model.Comment{
			Model:    gorm.Model{},
			Vid:      uint(vid),
			Uid:      u.ID,
			ParentID: 0,
			Content:  req.Content,
		}
		if err = commentDao.Create(comment); err != nil {
			util.LogrusObj.Error(err)
			code = e.ErrorDataBase
			return
		}
		//向redis存储评论id
		if err = cache.AddVideoComments(ctx, *req.Vid, strconv.Itoa(int(comment.ID))); err != nil {
			util.LogrusObj.Error(err)
			code = e.ErrorRedis
			return
		}
	} else {
		//回复评论
		//如果request没有指定vid，需要从数据库里查找vid
		if req.Vid == nil {
			cid, _ := strconv.ParseUint(*req.Cid, 10, 0)
			commentParent, err := commentDao.FindCommentByID(uint(cid))
			if err != nil {
				util.LogrusObj.Error(err)
				code = e.CommentNotFound
				return nil, code, err
			}
			vid := strconv.Itoa(int(commentParent.Vid))
			req.Vid = &vid
		}
		//向数据库存储数据
		vid, _ := strconv.Atoi(*req.Vid)
		parentID, _ := strconv.Atoi(*req.Cid)
		comment = &model.Comment{
			Model:    gorm.Model{},
			Vid:      uint(vid),
			Uid:      u.ID,
			ParentID: uint(parentID),
			Content:  req.Content,
		}
		if err = commentDao.Create(comment); err != nil {
			util.LogrusObj.Error(err)
			code = e.ErrorDataBase
			return
		}
		//在redis中 对主评论增加子评论id
		if err = cache.AddCommentChildren(ctx, *req.Cid, strconv.Itoa(int(comment.ID))); err != nil {
			util.LogrusObj.Debug(err)
			code = e.ErrorRedis
			return
		}
	}
	//构建返回结构
	resp = &interaction.CommentPublishResponse{Comment: types.BuildComment(ctx, comment)}
	return resp, e.SUCCESS, nil
}

func (s *InteractionService) CommentList(ctx context.Context, req *interaction.CommentListRequest) (resp *interaction.CommentListResponse, count, code int64, err error) {
	//判断参数
	if req.Vid == nil && req.Cid == nil {
		err = errors.New("缺少参数")
		util.LogrusObj.Info(err)
		code = e.ParamsMissing
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
	commentDao := dao.NewCommentDao(ctx)
	var list []string

	//视频评论列表
	if req.Vid != nil && req.Cid == nil {
		//从redis获取视频评论id
		list, err = cache.VideoCommentsList(ctx, *req.Vid)
		if err != nil {
			util.LogrusObj.Debug(err)
			code = e.ErrorRedis
			return
		}
	} else {
		//从redis获取评论的子评论id
		list, err = cache.CommentChildrenList(ctx, *req.Cid)
		if err != nil {
			util.LogrusObj.Debug(err)
			code = e.ErrorRedis
			return
		}
	}
	//从数据库获取评论
	res, count, err := commentDao.FindCommentsByID(list, *req.PageNum, *req.PageSize)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.CommentNotFound
		return
	}
	//构建返回结构
	resp = &interaction.CommentListResponse{Comments: types.BuildCommentList(ctx, res)}
	return
}

func (s *InteractionService) DeleteComment(ctx context.Context, req *interaction.DeleteCommentRequest) (code int64, err error) {
	commentDao := dao.NewCommentDao(ctx)

	if err = commentDao.DeleteComment(*req.Cid, *req.Vid); err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorDataBase
		return
	}
	return e.SUCCESS, nil
}
