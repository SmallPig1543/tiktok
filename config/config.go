package config

import (
	"github.com/spf13/viper"
	"os"
)

var Config Conf

type Conf struct {
	System   *System `mapstructure:"system"`
	Mysql    `mapstructure:"mysql"`
	Redis    `mapstructure:"redis"`
	Oss      `mapstructure:"oss"`
	Es       `mapstructure:"es"`
	Rabbitmq `mapstructure:"rabbitmq"`
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

type Es struct {
	UserName string `mapstructure:"esUserName"`
	Password string `mapstructure:"esPassword"`
	Host     string `mapstructure:"esHost"`
	Port     string `mapstructure:"esPort"`
}

type Rabbitmq struct {
	RabbitmqUser     string `mapstructure:"rabbitmqUserName"`
	RabbitmqPassword string `mapstructure:"rabbitmqPassword"`
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/config/local")
	viper.AddConfigPath(workDir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
}
