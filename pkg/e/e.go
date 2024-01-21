package e

const (
	SUCCESS int64 = 200
	ERROR   int64 = 500

	InvalidParams int64 = iota

	ErrorUserExist
	ErrorUserNotExist
	TokenGeneratedFail

	ErrorDataBase
	ErrorRedis
)
