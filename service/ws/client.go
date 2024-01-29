package ws

import (
	"encoding/json"
	"fmt"
	"github.com/hertz-contrib/websocket"
	"log"
	"tiktok/pkg/util"
)

type Client struct {
	Uid     string          `json:"uid"`
	SendID  string          `json:"send_id"`
	Message chan []byte     `json:"message"`
	Conn    *websocket.Conn `json:"conn"`
}

type Message struct {
	From    string `json:"from"`
	Content string `json:"content"`
}

func (c *Client) GetMsg() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Conn.Close()
	}()
	for {
		c.Conn.PongHandler()
		msg := &Message{}
		_, msgBuf, err := c.Conn.ReadMessage()
		if err != nil {
			util.LogrusObj.Debug(err)
			Manager.Unregister <- c
			_ = c.Conn.Close()
			return
		}
		err = json.Unmarshal(msgBuf, &msg)
		log.Println(msg.Content)
		Manager.Broadcast <- &Broadcast{
			Client:  c,
			Message: []byte(msg.Content),
		}

	}
}

func (c *Client) WriteMsg() {
	defer func() {
		_ = c.Conn.Close()
	}()
	for {
		select {
		//写之前先接收数据
		case message, ok := <-c.Message:
			if !ok {
				//没有消息就关掉通道
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			replyMsg := &Message{
				From:    c.Uid,
				Content: fmt.Sprintf("%s", string(message)),
			}
			util.LogrusObj.Debug(replyMsg.Content)
			msg, _ := json.Marshal(replyMsg)
			_ = c.Conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
