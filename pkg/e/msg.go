package e

var MsgFlags = map[int64]string{
	SUCCESS: "操作成功",
	ERROR:   "操作失败",

	InvalidParams: "请求参数有误",
	ParamsMissing: "缺少参数",

	ErrorDataBase: "数据库出错",
	ErrorRedis:    "redis出错",
	//token
	TokenGeneratedFail: "token生成失败",
	ErrorTokenTimeout:  "token超时",
	ParseTokenFailure:  "token解析出错",

	//user
	ErrorUserExist:         "用户已存在",
	ErrorUserNotExist:      "用户不存在",
	SetPasswordFail:        "密码设置失败",
	ErrorGetUserInfo:       "获取用户信息失败",
	VerifyOtpFailed:        "otp验证失败",
	UpdateTotpStatusFailed: "更新2FA状态失败",
	ErrorGenerateOTP:       "生成otp失败",
	ErrorPassword:          "密码错误",
	ErrorGetAvatar:         "获取头像失败",
	ErrorAvatarUpload:      "头像上传失败",
	InvalidAvatar:          "头像格式不符",
	ErrorOssUpload:         "云端上传出错",
	MFABindFail:            "MFA认证失败",
	OtpMissing:             "缺少otp",
	//video
	ErrorVideoOpen:     "视频打开失败",
	ErrorVideoUpload:   "视频上传失败",
	ErrorGetUrl:        "获取视频url失败",
	ErrorVideoNotExist: "视频不存在",
	InValidVideo:       "视频格式不符要求",
	InValidCover:       "封面不符要求",
	SaveVideoFail:      "视频保持本地失败",
	SaveCoverFail:      "封面保持本地失败",
	//comment
	ErrorComment:    "评论失败",
	CommentNotFound: "评论不存在",

	//document
	ErrorCreateDoc: "创建文档失败",

	//interaction
	ActionFails:  "操作失败",
	GetListFails: "获取列表失败",
}

// GetMsg 获取错误码对应的信息
func GetMsg(code int64) string {
	if msg, ok := MsgFlags[code]; ok {
		return msg
	}
	return MsgFlags[ERROR]
}
