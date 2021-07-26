package main

import (
	"flag"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/quicoment/qc-mq-processing-server/api"
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

	r := setupRouter()
	if err := r.Run(); err != nil {
		errors.Errorf("Fail gin engine start: %w", err)
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}))

	r.POST("/queues", api.CreateQueue)

	return r
}

func setupConsumer() {
	queueNames, err := common.GETALL_SET("QUEUE")

	if err != nil {
		errors.Errorf("Fail setup consumer: %w", err)
	}

	for _, queueName := range queueNames {
		var err error
		rabbitConfig := common.RabbitConfig{
			Schema:         "amqp",
			Username:       "username",
			Password:       "password",
			Host:           "127.0.0.1",
			Port:           "5672",
			VHost:          "",
			ConnectionName: "",
		}
		rabbit := common.NewRabbit(rabbitConfig)
		if err = rabbit.ConnectRabbit(); err != nil {
			log.Fatalf("unable to connect to rabbit: %w", err)
		}

		consumerConfig := common.ConsumerConfig{
			ExchangeName:  "name.test", // TODO: exchange name 설정 필요
			ExchangeType:  "direct",
			RoutingKey:    "create",
			QueueName:     queueName,
			ConsumerName:  queueName,
			ConsumerCount: 3,
			PrefetchCount: 1,
		}
		consumerConfig.Reconnect.MaxAttempt = 60
		consumerConfig.Reconnect.Interval = 1 * time.Second
		consumer := common.NewConsumer(consumerConfig, rabbit)
		if err = consumer.ConsumerStart(); err != nil {
			log.Fatalf("unable to start consumer: %w", err)
		}
	}
}
