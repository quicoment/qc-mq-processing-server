package common

import (
	"fmt"
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

	log.Println("[", id, "] Running ...")
	log.Println("[", id, "] Press CTRL+C to exit ...")
	for msg := range messages {
		log.Println("[", id, "] Consumed:", string(msg.Body))
		// TODO: parse message and redis update


		if err := msg.Ack(false); err != nil {
			// TODO: ack을 보내지 못했을 때
			log.Println("unable to acknowledge the message, dropped", err)
		}

		log.Println("[", id, "] Exiting ...")
	}
}

