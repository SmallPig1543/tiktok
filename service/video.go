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
	"tiktok/dal/db/dao"
	"tiktok/dal/db/model"
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
	//构建返回结构
	resp = &video.PublishResponse{Video: types.BuildVideo(ctx, v)}
	return
}

func (s *VideoService) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, code int64, err error, total int64) {
	//从数据库中查询
	videoDao := dao.NewVideoDao(ctx)
	id, _ := strconv.Atoi(req.UID)
	videos, total, err := videoDao.VideoList(uint(id), int(req.PageNum), int(req.PageSize))
	if err != nil {
		util.LogrusObj.Debug(err)
		code = e.ErrorDataBase
		return
	}
	//构建返回结构
	resp = &video.PublishListResponse{Videos: types.BuildVideoList(ctx, videos)}
	return
}
