package ws

import (
	"github.com/hertz-contrib/websocket"
	"sync"
	"tiktok/dal/rabbitmq"
	"tiktok/pkg/util"
)

var Manager = NewHub()

// Broadcast 源用户和广播内容
type Broadcast struct {
	Client  *Client
	Message []byte
}

// Hub 维护用户和信息
type Hub struct {
	//已经注册的用户
	Clients     map[string]*Client
	ClientsLock sync.RWMutex

	//存储的message，用于广播
	Broadcast chan *Broadcast

	//发出注册请求的用户
	Register chan *Client

	//发出取消注册的用户
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan *Broadcast),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
	}
}

func (h *Hub) Start() {
	for {
		select {
		case client := <-h.Register: //注册用户
			h.Clients[client.Uid] = client
			_ = client.Conn.WriteMessage(websocket.TextMessage, []byte("已连接至服务器"))
		case client := <-h.Unregister:
			_ = client.Conn.WriteMessage(websocket.TextMessage, []byte("连接已断开"))
			if _, ok := h.Clients[client.Uid]; ok {
				delete(h.Clients, client.Uid)
				close(client.Message)
			}
		case broadcast := <-h.Broadcast:
			message := broadcast.Message
			sendID := broadcast.Client.SendID
			//如果用户在注册范围内，那么直接将消息发送过去
			if client, ok := h.Clients[sendID]; ok {
				select {
				case client.Message <- message:
				default:
					close(client.Message)
					delete(Manager.Clients, client.Uid)
				}
			} else {
				//如果不在，则将消息发送到消息队列rabbitmq
				err := rabbitmq.PublishMsg(message, sendID)
				if err != nil {
					util.LogrusObj.Error("发送失败")
					return
				}
			}
		}
	}
}
