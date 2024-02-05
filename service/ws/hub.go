package ws

import (
	"github.com/hertz-contrib/websocket"
	"sync"
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
			//对所有注册过的用户进行广播
			sendID := broadcast.Client.SendID
			//循环找到sendID和id匹配
			//TODO 使用rabbitmq优化
			for id, conn := range h.Clients {
				if id != sendID {
					continue
				}
				select {
				case conn.Message <- message:
				default:
					close(conn.Message)
					delete(Manager.Clients, conn.Uid)
				}
			}
		}
	}
}
