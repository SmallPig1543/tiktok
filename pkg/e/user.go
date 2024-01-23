package e

const (
	SetPasswordFail int64 = 10000 + iota
	ErrorGetUserInfo
	OtpMissing
	VerifyOtpFailed
	UpdateTotpStatusFailed
	ErrorGenerateOTP
	ErrorPassword
	GenerateTokenFailure
	ParseTokenFailure
	ErrorTokenTimeout
	ErrorGetAvatar
	ErrorAvatarUpload
	ErrorOssUpload
	ErrorComment
	InvalidAvatar
	MFABindFail
)
