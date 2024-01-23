package es

import (
	"github.com/olivere/elastic/v7"
	"tiktok/config"
	"tiktok/dal/es/model"
)

var EsClient *elastic.Client

func LinkEs() {
	conf := config.Config.Es
	url := "http://" + conf.Host + ":" + conf.Port
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false), //是否开启集群嗅探
		elastic.SetBasicAuth("elastic", "123456"),
	)
	if err != nil {
		panic(err)
	}
	EsClient = client
	Init()
}
func Init() {
	_ = CreateIndex(model.Video{})
}
