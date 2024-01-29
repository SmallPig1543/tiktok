package types

import (
	"context"
	"strconv"
	"tiktok/biz/model/interaction"
	"tiktok/dal/cache"
	"tiktok/dal/db/model"
)

func BuildComment(ctx context.Context, comment *model.Comment) *interaction.Comment {
	var deletedAt string
	if comment.DeletedAt.Valid {
		deletedAt = comment.DeletedAt.Time.Format("2006-01-02 15:04:05")
	}
	likes := cache.CommentLikes(ctx, strconv.Itoa(int(comment.ID)))
	children, _ := cache.CommentChildren(ctx, strconv.Itoa(int(comment.ID)))
	return &interaction.Comment{
		ID:        strconv.Itoa(int(comment.ID)),
		UID:       strconv.Itoa(int(comment.Uid)),
		Vid:       strconv.Itoa(int(comment.Vid)),
		ParentID:  strconv.Itoa(int(comment.ParentID)),
		Likes:     likes,
		Children:  children,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: comment.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt: deletedAt,
	}
}

func BuildCommentList(ctx context.Context, comments []*model.Comment) []*interaction.Comment {
	resp := make([]*interaction.Comment, 0)
	for _, data := range comments {
		resp = append(resp, BuildComment(ctx, data))
	}
	return resp
}
