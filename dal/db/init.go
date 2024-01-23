package dao

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"tiktok/config"
	"tiktok/dal/db/model"
)

var DB *gorm.DB

func InitMySQL() {
	conf := config.Config.Mysql
	dsn := conf.UserName + ":" + conf.MysqlPassword + "@tcp(" + conf.DbHost + ":" + conf.DbPort + ")/" + conf.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(&model.User{})
	_ = db.AutoMigrate(&model.Video{})
	DB = db
}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := DB
	return db.WithContext(ctx)
}
