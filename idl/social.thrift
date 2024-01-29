namespace go social

struct User {
    1: string id,
    2: string username
    3: string avatar_url
}

struct FollowRequest {
    1: required string uid
    2: required string action_type
}

struct FollowResponse {
}

struct FollowListRequest {
    1: required string uid
    2: optional i64 page_num
    3: optional i64 page_size
}

struct FollowListResponse {
    1: required list<User> users
}

struct FansListRequest {
    1: required string uid
    2: optional i64 page_num
    3: optional i64 page_size
}

struct FansListResponse {
    1:required list<User> users
}

struct FriendsListRequest {
    1: optional i64 page_num
    2: optional i64 page_size
}

struct FriendsListResponse {
    1: required list<User> users
}

service SocialService {
    FollowResponse Follow(1: FollowRequest req) (api.post="tiktok/relation/action")
    FollowListResponse FollowList(1: FollowListRequest req) (api.get="tiktok/following/list")
    FansListResponse FansList(1: FansListRequest req) (api.get="tiktok/follower/list")
    FriendsListResponse FriendsList(1: FriendsListRequest req) (api.get="tiktok/friends/list")
}