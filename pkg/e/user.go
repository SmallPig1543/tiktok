package e

const (
	SetPasswordFail int64 = 10000 + iota
	ErrorGetUserInfo
	VerifyOtpFailed
	UpdateTotpStatusFailed
	ErrorGenerateOTP
	ErrorPassword
	GenerateTokenFailure
	ParseTokenFailure
	ErrorTokenTimeout
	ErrorGetAvatar
	ErrorAvatarUpload
	ErrorComment
)
