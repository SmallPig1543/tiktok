package errno

var (
	Success    = NewErrNo(SuccessCode, "Success")
	ServiceErr = NewErrNo(ServiceErrCode, "Service is unable to start successfully")
	ParamErr   = NewErrNo(ParamErrCode, "Wrong Parameter has been given")
	// user
	UserAlreadyExistErr    = NewErrNo(UserAlreadyExistErrCode, "User already exists")
	UserNotExist           = NewErrNo(UserNotExistErrCode, "User doesn't exists")
	AuthorizationFailedErr = NewErrNo(AuthorizationFailedErrCode, "Authorization failed")

	//utils
	UploadErr = NewErrNo(UploadErrCode, "Upload failed")
)
