// Code generated by hertz generator.

package user

import (
	"github.com/cloudwego/hertz/pkg/app"
	"tiktok/biz/middleware"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _tiktokMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getinfoMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		middleware.JWT(),
	}
}

func _loginMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _registerMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _authMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		middleware.JWT(),
	}
}

func _mfaMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _mfabindMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getmfaqrcodeMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _avatarMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _avataruploadMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		middleware.JWT(),
	}
}
