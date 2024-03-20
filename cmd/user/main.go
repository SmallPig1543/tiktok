package main

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/spf13/viper"
	"tiktok/cmd/user/dal/db"
	"tiktok/config"
	user "tiktok/kitex_gen/user/userservice"
)

func init() {
	//初始化配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&config.Config); err != nil {
		panic(err)
	}
	//
	db.InitMySQL()

	klog.SetLevel(klog.LevelDebug)
}

func main() {
	conf := config.Config
	r, err := etcd.NewEtcdRegistry([]string{conf.EtcdHost + ":" + conf.EtcdPort})
	if err != nil {
		klog.Fatal(err)
	}

	svr := user.NewServer(new(UserServiceImpl), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "user"}), server.WithRegistry(r))

	err = svr.Run()

	if err != nil {
		panic(err)
	}
}
