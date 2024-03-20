package service

import "sync"

var UserServiceOnce sync.Once

var UserServiceIns *UserService

type UserService struct {
}

func GetUserService() *UserService {
	UserServiceOnce.Do(func() {
		UserServiceIns = &UserService{}
	})
	return UserServiceIns
}
