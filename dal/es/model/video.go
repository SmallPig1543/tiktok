package model

type Video struct {
	Vid         uint   `json:"vid"`
	Uid         uint   `json:"uid"`
	UserName    string `json:"user_name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreateAt    string `json:"create_at"`
}

func (v Video) Index() string {
	return "video_index"
}

func (v Video) Mapping() string {
	return `{
  "mappings": {
    "properties": {
      "uid": {
        "type": "integer"
      },
      "vid": {
        "type": "integer"
      },
      "title": {
        "type": "text",
        "analyzer": "ik_smart"
      },
      "user_name": {
		    "type": "keyword"
      },  
      "description": {
        "type": "text",
        "analyzer": "ik_smart"
      },
      "created_at": {
        "type": "date",
        "null_value": "null",
        "format": "[yyyy-MM-dd HH:mm:ss]"
      }
    }
  }
}`
}
