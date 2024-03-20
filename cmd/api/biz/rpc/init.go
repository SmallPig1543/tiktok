package rpc

import "tiktok/kitex_gen/user/userservice"

var (
	userClient userservice.Client
)

func Init() {
	InitUserRPC()
}
