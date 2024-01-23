package types

import (
	"context"
	"strconv"
	"tiktok/biz/model/video"
	"tiktok/dal/cache"
	"tiktok/dal/db/model"
	"tiktok/pkg/util"
)

func BuildVideo(ctx context.Context, v *model.Video) *video.Video {
	var deletedAt, coverURL, videoURL string
	if v.DeletedAt.Valid {
		deletedAt = v.DeletedAt.Time.Format("2006-01-02 15:04:05")
	}
	if v.CoverKey != "" {
		coverURL, _ = util.GetURL(v.CoverKey)
	}
	if v.VideoKey != "" {
		videoURL, _ = util.GetURL(v.VideoKey)
	}
	//从redis获取数据
	likes := cache.VideoLikes(ctx, v.ID)
	views := cache.VideoViews(ctx, v.ID)
	comments := cache.VideoComments(ctx, v.ID)
	return &video.Video{
		ID:          strconv.Itoa(int(v.ID)),
		UID:         strconv.Itoa(int(v.Uid)),
		URL:         videoURL,
		CoverURL:    coverURL,
		Title:       v.Title,
		Description: v.Description,
		Likes:       likes,
		Views:       views,
		Comments:    comments,
		CreatedAt:   v.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   v.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:   deletedAt,
	}
}

func BuildVideoList(ctx context.Context, videos []*model.Video) []*video.Video {
	resp := make([]*video.Video, 0)
	for _, data := range videos {
		resp = append(resp, BuildVideo(ctx, data))
	}
	return resp
}
