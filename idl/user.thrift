namespace go user


struct User {
    1: string id
    2: string name
    3: string avatar_url
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
}

//登录
struct LoginRequest {
    1: required string username
    2: required string password
    3: optional string otp
}

struct LoginResponse {
    1: required User user
    2: required string refresh_token
    3: required string access_token
}

//用户信息获取
struct InfoRequest {
    1: required string id
}

struct InfoResponse {
    1: required User user
}

//头像上传
struct AvatarUploadRequest {
    1: required binary data
}

struct AvatarUploadResponse {
    1: required User user
}

//获取MFAqrcode
struct GetMFAqrcodeRequest {
}

struct GetMFAqrcodeResponse {
    1: required string secret
    2: required string qrcode
}

//MFA认证
struct MFABindRequest {
    1: required string code
    2: required string secret
}

struct MFABindResponse {
}

service UserService {
    RegisterResponse Register(1: RegisterRequest req) (api.post="tiktok/user/register")
    LoginResponse Login(1: LoginRequest req) (api.post="tiktok/user/login")
    InfoResponse GetInfo(1 :InfoRequest req) (api.get="tiktok/user/info")
    AvatarUploadResponse AvatarUpload(1: AvatarUploadRequest req) (api.put="tiktok/user/avatar/upload")
    GetMFAqrcodeResponse GetMFAqrcode(1: GetMFAqrcodeRequest req) (api.get="tiktok/auth/mfa/qrcode")
    MFABindResponse MFABind(1: MFABindRequest req) (api.post="tiktok/auth/mfa/bind")
}