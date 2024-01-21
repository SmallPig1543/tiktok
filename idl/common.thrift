namespace go common

struct BaseResp {
    1:i64 code
    2:string msg
}

service CommonService {
    BaseResp Ping() (api.get="tiktok/ping")
}