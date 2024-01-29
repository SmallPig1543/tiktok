package document

import (
	"context"
	"tiktok/dal/es"
	"tiktok/dal/es/model"
)

func CreateDocument(model model.Video, id string, ctx context.Context) (err error) {
	_, err = es.EsClient.Index().Index(model.Index()).Id(id).BodyJson(&model).Do(ctx)
	return
}

func UpdateDocument(indexName string, key string, value string, id string, ctx context.Context) (err error) {
	_, err = es.EsClient.Update().Index(indexName).Id(id).Doc(map[string]any{
		key: value,
	}).Do(ctx)
	return
}
