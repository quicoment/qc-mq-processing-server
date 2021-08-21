package common

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/quicoment/qc-mq-processing-server/config"
	"github.com/quicoment/qc-mq-processing-server/domain"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

type ConsumerConfigs struct {
	Configs []ConsumerConfig `yaml:"ConsumerConfigs"`
}

type ConsumerConfig struct {
	QueueName     string `yaml:"QueueName"`
	ConsumerName  string `yaml:"ConsumerName"`
	ConsumerCount int    `yaml:"ConsumerCount"`
	PrefetchCount int    `yaml:"PrefetchCount"`
	Reconnect     struct {
		MaxAttempt int
		Interval   time.Duration
	} `yaml:"Reconnect"`
}

type Consumer struct {
	config ConsumerConfig
	Rabbit *Rabbit
}

func NewConsumer(config ConsumerConfig, rabbit *Rabbit) *Consumer {
	return &Consumer{
		config: config,
		Rabbit: rabbit,
	}
}

func (c *Consumer) ConsumerStart() error {
	con, err := c.Rabbit.CheckRabbitConnection()
	if err != nil {
		log.Fatal("unable to get rabbit connection: %w", err)
		return err
	}
	go c.closedConnectionListener(con.NotifyClose(make(chan *amqp.Error)))

	chn, err := con.Channel()
	if err != nil {
		log.Fatal("unable to get rabbit channel: %w", err)
		return err
	}

	if err := chn.Qos(c.config.PrefetchCount, 0, false); err != nil {
		return err
	}

	for i := 1; i <= c.config.ConsumerCount; i++ {
		id := i
		go c.consume(chn, id)
	}

	return nil
}

func (c *Consumer) closedConnectionListener(closed <-chan *amqp.Error) {
	log.Println("INFO: Watching closed connection")

	err := <-closed
	if err != nil {
		log.Println("INFO: Closed connection:", err.Error())

		var i int

		for i = 0; i < c.config.Reconnect.MaxAttempt; i++ {
			log.Println("INFO: Attempting to reconnect")

			if err := c.Rabbit.ConnectRabbit(); err == nil {
				log.Println("INFO: Reconnected")

				if err := c.ConsumerStart(); err == nil {
					break
				}
			}

			time.Sleep(c.config.Reconnect.Interval)
		}

		if i == c.config.Reconnect.MaxAttempt {
			log.Println("CRITICAL: Giving up reconnecting")

			return
		}
	} else {
		log.Println("INFO: Connection closed normally, will not reconnect")
		os.Exit(0)
	}
}

func (c *Consumer) consume(channel *amqp.Channel, id int) {
	messages, err := channel.Consume(
		c.config.QueueName,
		fmt.Sprintf("%s (%d/%d)", c.config.ConsumerName, id, c.config.ConsumerCount),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(fmt.Sprintf("CRITICAL: Unable to start consumer (%d/%d): %w", id, c.config.ConsumerCount, err))

		return
	}

	log.Printf("[%d] Running ...", id)
	log.Printf("[%d] Press CTRL+C to exit ...", id)
	for msg := range messages {
		log.Printf("[%d] Consumed:", id, string(msg.Body))

		var body map[string]interface{}
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			log.Printf("[%d] Consumed but cannot parse to json: %s", id, string(msg.Body))
		}

		if err := processMessage(fmt.Sprintf("%v", msg.Body), fmt.Sprintf("%v", body["messageType"])); err != nil {
			_ = msg.Reject(true)
			log.Printf("[%d] Consumed but error with process message: %w", id, err)
		} else if err := msg.Ack(false); err != nil { // ACK 을 보내지 못했을 때
			_ = msg.Reject(true)
			log.Printf("unable to acknowledge the message, dropped: %w", err)
		}

		log.Printf("[%d] Exiting ...", id)
	}
}

func processMessage(message string, messageType string) error {
	switch messageType {
	case "register":
		var request domain.CommentCreateRequest
		if err := json.Unmarshal([]byte(message), &request); err != nil {
			errors.Errorf("Cannot parse to create request json: %w", err)
		}
		request.ID = config.GetId()
		comment := domain.NewComment(request)
		return createComment(comment)

	case "like":
		var commentLikeRequest domain.CommentLikeRequest
		if err := json.Unmarshal([]byte(message), &commentLikeRequest); err != nil {
			errors.Errorf("Cannot parse to create request json: %w", err)
		}
		return likeComment(commentLikeRequest.UserId, commentLikeRequest.PostId, commentLikeRequest.CommentId)

	default:
		return errors.Errorf("wrong message type: %s", messageType)
	}
}
