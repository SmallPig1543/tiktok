package ws

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/websocket"
	"strconv"
	"tiktok/dal/rabbitmq"
	"tiktok/pkg/ctl"
	"tiktok/pkg/util"
)

var upgrader = websocket.HertzUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(ctx *app.RequestContext) bool {
		return true
	},
}

func ServerWs(ctx context.Context, c *app.RequestContext) {
	//升级为websocket
	u, err := ctl.GetUserInfo(ctx)
	uid := strconv.Itoa(int(u.ID))
	toUid := c.Query("to_uid")
	err = upgrader.Upgrade(c, func(conn *websocket.Conn) {
		client := &Client{
			Uid:     uid,
			SendID:  toUid,
			Message: make(chan []byte, 256),
			Conn:    conn,
		}
		Manager.Register <- client
		list, err := rabbitmq.Consume(uid)
		if err != nil {
			return
		}
		for _, v := range list {
			client.Message <- v
		}
		go client.WriteMsg()
		client.GetMsg()
	})
	if err != nil {
		util.LogrusObj.Debug(err)
	}
}
