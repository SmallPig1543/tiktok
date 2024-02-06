package rabbitmq

import (
	"fmt"
	"tiktok/pkg/util"
)

func Consume(to string) (list [][]byte, err error) {
	//声明通道
	ch, err := RabbitmqConn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()
	//声明交换机
	err = ch.ExchangeDeclare(
		"direct_exchange", // Exchange 名称
		"direct",          // Exchange 类型
		true,              // 持久化
		false,             // 不自动删除
		false,             // 不等待服务器响应
		false,
		nil, // 不设置额外参数
	)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}
	// 声明队列
	queue, err := ch.QueueDeclare(
		to,    // 队列名称
		false, // 不持久化
		false, // 不自动删除
		false, // 不独占
		false, // 不等待服务器响应
		nil,   // 不设置额外参数
	)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}

	// 交换机和队列绑定
	err = ch.QueueBind(
		queue.Name,        // 队列名称
		to,                // 路由键，用于绑定 Exchange 和队列
		"direct_exchange", // Exchange 名称
		false,             // 不等待服务器响应
		nil,               // 不设置额外参数
	)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}
	list = make([][]byte, 0)
	for {
		msg, ok, err := ch.Get(to, true)
		if err != nil {
			util.LogrusObj.Error(err)
			return nil, err
		}
		if !ok {
			break
		}
		list = append(list, msg.Body)
	}
	return
}
