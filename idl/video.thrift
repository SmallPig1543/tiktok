namespace go video

struct Video {
    1: string id
    2: string uid
    3: string url
    4: string cover_url
    5: string title
    6: string description
    7: i64 likes
    8: i64 views
    9: i64 comments
    10: string created_at
    11: string updated_at
    12: string deleted_at
}

struct FeedRequest {
    1: required i64 timeStamp
}

struct FeedResponse {
    1: required list<Video> videos
}

struct PublishRequest {
    //1: required binary data
    1: required string title
    2: required string description
}

struct PublishResponse {
    1: required Video video
}

struct PublishListRequest {
    1: required string uid
    2: required i64 page_num
    3: required i64 page_size
}

struct PublishListResponse {
    1: required list<Video> videos
}

struct PopularListRequest {
    1: required i64 page_num
    2: required i64 page_size
}

struct PopularListResponse {
    1: required list<Video> videos
}

struct SearchRequest {
    1: required string keyword
    2: required i64 page_num
    3: required i64 page_size
    4: optional i64 from_date
    5: optional i64 to_date
    6: optional string username
}

struct SearchResponse {
    1: required list<Video> videos
}

service VideoService {
    FeedResponse Feed(1: FeedRequest req) (api.get="tiktok/video/feed")
    PublishResponse Publish(1: PublishRequest req) (api.post="tiktok/video/publish")
    PublishListResponse PublishList(1: PublishListRequest req) (api.get="tiktok/video/list")
    PopularListResponse PopularList(1: PopularListRequest req) (api.get="tiktok/video/popular")
    SearchResponse Search(1: SearchRequest req) (api.post="tiktok/video/search")
}