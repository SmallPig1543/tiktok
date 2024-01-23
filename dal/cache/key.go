package cache

import "fmt"

func VideoViewKey(vid uint) string {
	return fmt.Sprintf("video:view:%d", vid)
}

func VideoInfoKey(vid uint) string {
	return fmt.Sprintf("video:info:%d", vid)
}

func VideoLikeKey(vid uint) string {
	return fmt.Sprintf("video:like:%d", vid)
}

func VideoCommentKey(vid uint) string {
	return fmt.Sprintf("video:comment:%d", vid)
}

func ChildrenComments(cid uint) string {
	return fmt.Sprintf("comment:childrenCount:%d", cid)
}
