package rabbitmq

import (
	"github.com/streadway/amqp"
	"tiktok/config"
	"tiktok/pkg/util"
)

var RabbitmqConn *amqp.Connection

func LinkRabbitmq() {
	conf := config.Config.Rabbitmq
	conn, err := amqp.Dial("amqp://" + conf.RabbitmqUser + ":" + conf.RabbitmqPassword + "@" + conf.RabbitmqHost + ":" + conf.RabbitmqPort + "/")
	if err != nil {
		panic(err)
	}
	RabbitmqConn = conn
	util.LogrusObj.Info("connecting to rabbitmq successes")
}
