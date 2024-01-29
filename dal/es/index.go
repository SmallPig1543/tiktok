package es

import (
	"context"
	"tiktok/dal/es/model"
)

func CreateIndex(model model.Video) (err error) {
	_, err = EsClient.CreateIndex(model.Index()).BodyString(model.Mapping()).Do(context.Background())
	return
}

func IndexExists(indexName string) bool {
	exists, _ := EsClient.IndexExists(indexName).Do(context.Background())
	return exists
}

func DeleteIndex(indexName string) (err error) {
	_, err = EsClient.DeleteIndex(indexName).Do(context.Background())
	return
}
