package config

import (
	"github.com/spf13/viper"
	"os"
)

var Config Conf

type Conf struct {
	System *System `mapstructure:"system"`
	Mysql  `mapstructure:"mysql"`
	Redis  `mapstructure:"redis"`
	Oss    `mapstructure:"oss"`
	Email  `mapstructure:"email"`
	Es     `mapstructure:"es"`
}

type System struct {
	Domain string `mapstructure:"domain"`
	Host   string `mapstructure:"host"`
	Port   string `mapstructure:"port"`
}

type Mysql struct {
	UserName      string `mapstructure:"userName"`
	MysqlPassword string `mapstructure:"mysqlPassword"`
	DbName        string `mapstructure:"dbName"`
	DbHost        string `mapstructure:"dbHost"`
	DbPort        string `mapstructure:"dbPort"`
}

type Redis struct {
	RedisHost     string `mapstructure:"redisHost"`
	RedisPort     string `mapstructure:"redisPort"`
	RedisPassword string `mapstructure:"redisPassword"`
	RedisDbName   int    `mapstructure:"redisDbName"`
}

type Oss struct {
	OssEndPoint        string `mapstructure:"OSS_END_POINT"`
	OssAccessKeyId     string `mapstructure:"OSS_ACCESS_KEY_ID"`
	OssAccessKeySecret string `mapstructure:"OSS_ACCESS_KEY_SECRET"`
	OssBucket          string `mapstructure:"OSS_BUCKET"`
}

type Email struct {
	Sender   string `mapstructure:"sender"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Addr     string `mapstructure:"addr"`
}

type Es struct {
	UserName string `mapstructure:"esUserName"`
	Password string `mapstructure:"esPassword"`
	Host     string `mapstructure:"esHost"`
	Port     string `mapstructure:"esPort"`
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	filePath := workDir + "/config/local"
	viper.AddConfigPath(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
}
