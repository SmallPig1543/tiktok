package document

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"strconv"
	"tiktok/biz/model/video"
	"tiktok/dal/es"
	"tiktok/dal/es/model"
	"tiktok/pkg/util"
	"time"
)

func SearchVideo(ctx context.Context, req *video.SearchRequest) (list []string, err error) {
	query := elastic.NewBoolQuery()
	//匹配关键词
	query = query.Should(elastic.NewMatchQuery("title", req.Keyword))
	query = query.Should(elastic.NewMatchQuery("description", req.Keyword))

	if req.Username != nil {
		query = query.Must(elastic.NewTermsQuery("user_name", *req.Username))
	}
	//时间范围查找
	if req.FromDate != nil || req.ToDate != nil {
		var startTime, endTime string
		if req.FromDate != nil {
			startTime = time.Unix(0, *req.FromDate).Format("2006-01-02 15:04:05")
		}
		if req.ToDate != nil {
			endTime = time.Unix(0, *req.FromDate).Format("2006-01-02 15:04:05")
		}
		if req.FromDate != nil && req.ToDate != nil {
			query = query.Must(elastic.NewRangeQuery("create_at").Gte(startTime).Lte(endTime))
		} else if req.FromDate != nil {
			query = query.Must(elastic.NewRangeQuery("create_at").Gte(startTime))
		} else {
			query = query.Must(elastic.NewRangeQuery("create_at").Lte(endTime))
		}
	}
	start := (*req.PageNum - 1) * (*req.PageSize)
	if start < 0 {
		start = 0
	}
	res, err := es.EsClient.Search(model.Video{}.Index()).Query(query).From(int(start)).Size(int(*req.PageSize)).Do(ctx)
	if err != nil {
		util.LogrusObj.Debug(err)
		return
	}
	var v model.Video
	list = make([]string, 0)
	for _, data := range res.Hits.Hits {
		_ = json.Unmarshal(data.Source, &v)
		list = append(list, strconv.Itoa(int(v.Vid)))
	}
	return list, nil
}
