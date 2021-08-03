package main

import (
	"flag"
	"github.com/quicoment/qc-mq-processing-server/common"
	"log"
	"time"
)

var (
	redisServer = flag.String("127.0.0.1", ":6379", "redis-connect-host")
)

func main() {
	common.InitRedisPool(*redisServer)
	setupConsumer()
	select {}
}

func setupConsumer() {
	rabbitConfig := common.RabbitConfig{
		Schema:         "amqp",
		Username:       "username",
		Password:       "password",
		Host:           "127.0.0.1",
		Port:           "5672",
		VHost:          "",
		ConnectionName: "amqp",
	}

	consumerConfig := common.ConsumerConfig{
		ConsumerName:  "q.quicoment.name",
		ConsumerCount: 100,
		PrefetchCount: 100, // 메세지까지 한번에 Listener의 메모리에 Push 한 뒤 Consumer 가 메모리에서 하나씩 메세지를 꺼내서 처리
	}

	rabbit := common.NewRabbit(rabbitConfig)

	var err error
	if err = rabbit.ConnectRabbit(); err != nil {
		log.Fatalf("unable to connect to rabbit: %w", err)
	}

	consumerConfig.Reconnect.MaxAttempt = 60
	consumerConfig.Reconnect.Interval = 1 * time.Second
	consumer := common.NewConsumer(consumerConfig, rabbit)

	if err = consumer.ConsumerStart(); err != nil {
		log.Fatalf("unable to start consumer: %w", err)
	}
}
