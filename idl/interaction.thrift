namespace go interaction

include "video.thrift"

struct Comment {
    1: string id
    2: string uid
    3: string vid
    4: string parent_id
    5: i64 likes
    6: i64 children
    7: string content
    8: string created_at
    9: string updated_at
    10: string deleted_at
}
//点赞操作
struct LikeRequest {
    1: optional string vid
    2: optional string cid
    3: required string action_type
}

struct LikeResponse {
}

//指定用户点赞的视频列表
struct LikeListRequest {
    1: required string uid
    2: optional i64 page_num
    3: optional i64 page_size
}

struct LikeListResponse {
    1: required list<video.Video> videos
}

//发布评论
struct CommentPublishRequest {
    1: optional string vid
    2: optional string cid
    3: required string content
}

struct CommentPublishResponse {
    1: required Comment comment
}

//评论列表
struct CommentListRequest {
    1: optional string vid
    2: optional string cid
    3: optional i64 page_num
    4: optional i64 page_size
}

struct CommentListResponse {
    1: optional list<video.Video> videos
    2: optional list<Comment> comments
}

//删除评论
struct DeleteCommentRequest {
    1: optional string vid
    2: optional string cid
}

struct DeleteCommentResponse {
}

service InteractionService {
    LikeResponse Like(1: LikeRequest req) (api.post="tiktok/like/action")
    LikeListResponse LikeList(1: LikeListRequest req) (api.get="tiktok/like/list")
    CommentPublishResponse CommentPublish(1: CommentListRequest req) (api.post="tiktok/comment/publish")
    CommentListResponse CommentList(1: CommentListRequest req) (api.get="tiktok/comment/list")
    DeleteCommentResponse DeleteComment(1: DeleteCommentRequest req) (api.delete="tiktok/comment/delete")
}
