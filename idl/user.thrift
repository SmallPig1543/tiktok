namespace go user

struct BaseResp {
    1: i64 code
    2: string msg
}

struct User {
    1: string id
    2: string name
    3: string avatar
    4: string created_at
    5: string updated_at
    6: string deleted_at
}

//注册
struct RegisterRequest {
    1: required string username
    2: required string password
}

struct RegisterResponse {
    1: BaseResp baseResp
}

//登录
struct LoginRequest {
    1: string username
    2: string password
    3: optional string otp
}

struct LoginResponse {
    1:BaseResp baseResp
    2: User user
    3: string access_token
    4: string refresh_token
}

//用户信息获取
struct InfoRequest {
    1: string uid
}

struct InfoResponse {
    1: BaseResp baseResp
    2: User user
}

//头像上传
struct AvatarUploadRequest {
    1: binary data
    2: string uid
}

struct AvatarUploadResponse {
    1: BaseResp baseResp
    2: User user
}

//获取MFAqrcode
struct GetMFAqrcodeRequest {
    1: string uid
    2: string username
}

struct GetMFAqrcodeResponse {
    1: BaseResp baseResp
    2: string secret
    3: string qrcode
}

//MFA认证
struct MFABindRequest {
    1: string code
    2: string secret
    3: string uid
}

struct MFABindResponse {
    1:BaseResp baseResp
}

service UserService {
    RegisterResponse Register(1: RegisterRequest req)
    LoginResponse Login(1: LoginRequest req)
    InfoResponse GetInfo(1 :InfoRequest req)
    AvatarUploadResponse AvatarUpload(1: AvatarUploadRequest req)
    GetMFAqrcodeResponse GetMFAqrcode(1: GetMFAqrcodeRequest req)
    MFABindResponse MFABind(1: MFABindRequest req)
}