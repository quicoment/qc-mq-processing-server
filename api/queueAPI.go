package api

import (
	"github.com/apsdehal/go-logger"
	"github.com/gin-gonic/gin"
	"github.com/quicoment/qc-mq-processing-server/common"
	"github.com/quicoment/qc-mq-processing-server/domain"
	"net/http"
	"time"
)

var (
	log, _ = logger.New("api", 1)
)

func CreateQueue(c *gin.Context) {
	var request domain.QueueCreateRequest

	var err error
	if err = c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body type error"})
		return
	}

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
		QueueName:     request.QueueName,
		ConsumerName:  request.QueueName,
		ConsumerCount: 3,
		PrefetchCount: 1,
	}
	consumerConfig.Reconnect.MaxAttempt = 60
	consumerConfig.Reconnect.Interval = 1 * time.Second
	consumer := common.NewConsumer(consumerConfig, rabbit)
	if err = consumer.ConsumerStart(); err != nil {
		log.Fatalf("unable to start consumer: %w", err)
	}

	if err = common.INSERT_SET("QUEUE", request.QueueName); err != nil {
		log.Fatalf("unable to insert queue name in redis: %w", err)
	}

	c.JSON(http.StatusCreated, gin.H{})
}
