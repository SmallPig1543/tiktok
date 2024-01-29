package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
	"path"
	"strconv"
	"strings"
	"sync"
	"tiktok/biz/model/video"
	"tiktok/dal/cache"
	"tiktok/dal/db/dao"
	"tiktok/dal/db/model"
	"tiktok/dal/es/document"
	esmodel "tiktok/dal/es/model"
	"tiktok/pkg/ctl"
	"tiktok/pkg/e"
	"tiktok/pkg/util"
	"tiktok/types"
	"time"
)

var VideoServiceOnce sync.Once
var VideoServiceIns *VideoService

type VideoService struct {
}

func GetVideoService() *VideoService {
	VideoServiceOnce.Do(func() {
		VideoServiceIns = &VideoService{}
	})
	return VideoServiceIns
}

func (s *VideoService) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, code int64, err error, count int64) {
	videoDao := dao.NewVideoDao(ctx)
	targetTime := time.UnixMilli(req.TimeStamp)
	list, count, err := videoDao.VideoFeed(targetTime)
	resp = &video.FeedResponse{Videos: types.BuildVideoList(ctx, list)}
	return
}

func (s *VideoService) Publish(ctx context.Context, req *video.PublishRequest, videoData, coverData *multipart.FileHeader) (resp *video.PublishResponse, code int64, err error) {
	//先判断文件类型
	if ok := util.IsVideo(videoData); !ok {
		err = errors.New("无效视频格式")
		code = e.InValidVideo
		return
	}
	if ok := util.IsImage(coverData); !ok {
		err = errors.New("封面不符要求")
		code = e.InValidCover
		return
	}
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorGetUserInfo
		return
	}
	//将视频下载到本地，方便上传到oss
	ext := strings.ToLower(path.Ext(videoData.Filename))
	storePath := "./static/video/" + u.UserName + uuid.Must(uuid.NewRandom()).String() + ext
	if err = util.SaveFile(videoData, storePath); err != nil {
		util.LogrusObj.Debug(err)
		code = e.SaveVideoFail
		return
	}
	//上传oss，获取key
	videoKey, err := util.VideoUpload(storePath)
	if err != nil {
		util.LogrusObj.Debug(err)
		code = e.ErrorOssUpload
		return
	}
	//相同的原理保持封面
	ext = strings.ToLower(path.Ext(coverData.Filename))
	storePath = "./static/cover/" + u.UserName + uuid.Must(uuid.NewRandom()).String() + ext
	if err = util.SaveFile(coverData, storePath); err != nil {
		util.LogrusObj.Debug(err)
		code = e.SaveCoverFail
		return
	}
	coverKey, err := util.CoverUpload(storePath)
	if err != nil {
		util.LogrusObj.Debug(err)
		code = e.ErrorOssUpload
		return
	}
	//更新数据库
	v := &model.Video{
		Model:       gorm.Model{},
		Uid:         u.ID,
		UserName:    u.UserName,
		Title:       req.Title,
		Description: req.Description,
		CoverKey:    coverKey,
		VideoKey:    videoKey,
	}
	videoDao := dao.NewVideoDao(ctx)
	if err = videoDao.CreateVideo(v); err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorDataBase
		return
	}
	// 将数据放入es
	esModel := esmodel.Video{
		Vid:         v.ID,
		Uid:         v.Uid,
		UserName:    v.UserName,
		Title:       v.Title,
		Description: v.Description,
		CreateAt:    v.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if err = document.CreateDocument(esModel, strconv.Itoa(int(v.ID)), ctx); err != nil {
		util.LogrusObj.Debug(err)
		code = e.ErrorCreateDoc
		return
	}
	//构建返回结构
	resp = &video.PublishResponse{Video: types.BuildVideo(ctx, v)}
	return
}

func (s *VideoService) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, count, code int64, err error) {
	//判断一下参数
	if req.PageSize == nil {
		req.PageSize = new(int64)
		*req.PageSize = 10
	}
	if req.PageNum == nil {
		req.PageNum = new(int64)
		*req.PageNum = 0
	}
	//从数据库中查询
	videoDao := dao.NewVideoDao(ctx)
	id, _ := strconv.Atoi(req.UID)
	videos, count, err := videoDao.VideoList(uint(id), int(*req.PageNum), int(*req.PageSize))
	if err != nil {
		util.LogrusObj.Debug(err)
		code = e.ErrorDataBase
		return
	}
	//构建返回结构
	resp = &video.PublishListResponse{Videos: types.BuildVideoList(ctx, videos)}
	return
}

func (s *VideoService) PopularList(ctx context.Context, req *video.PopularListRequest) (resp *video.PopularListResponse, count, code int64, err error) {
	//判断一下参数
	if req.PageSize == nil {
		req.PageSize = new(int64)
		*req.PageSize = 10
	}
	if req.PageNum == nil {
		req.PageNum = new(int64)
		*req.PageNum = 0
	}

	//从redis获取vid
	start := (*req.PageNum - 1) * (*req.PageSize)
	end := start + (*req.PageSize) - 1
	list, err := cache.RedisClient.ZRevRange(ctx, "Rank", start, end).Result()
	if err != nil {
		util.LogrusObj.Error(err)
		code = e.ErrorRedis
		return
	}

	//从数据库获取数据
	var data *model.Video
	videoDao := dao.NewVideoDao(ctx)
	res := make([]*model.Video, 0)
	for _, v := range list {
		data, err = videoDao.FindVideoByVid(v)
		if err != nil {
			util.LogrusObj.Error(err)
			code = e.ErrorDataBase
			return
		}
		res = append(res, data)
	}
	resp = &video.PopularListResponse{Videos: types.BuildVideoList(ctx, res)}
	return
}

func (s *VideoService) Search(ctx context.Context, req *video.SearchRequest) (resp *video.SearchResponse, count, code int64, err error) {
	//处理请求参数
	if req.PageSize == nil {
		req.PageSize = new(int64)
		*req.PageSize = 10
	}
	if req.PageNum == nil {
		req.PageNum = new(int64)
		*req.PageNum = 0
	}
	list, err := document.SearchVideo(ctx, req)
	if err != nil {
		code = e.ErrorSearchDoc
		return
	}
	//从数据库获取数据
	res := make([]*model.Video, 0)
	var data *model.Video
	videoDao := dao.NewVideoDao(ctx)
	for _, s := range list {
		data, err = videoDao.FindVideoByVid(s)
		if err != nil {
			util.LogrusObj.Error(err)
			code = e.ErrorDataBase
			return
		}
		res = append(res, data)
	}
	count = int64(len(res))
	resp = &video.SearchResponse{Videos: types.BuildVideoList(ctx, res)}
	return
}
