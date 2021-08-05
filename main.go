package main

import (
	"fmt"
	"github.com/quicoment/qc-mq-processing-server/common"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	redisFile          = "/values/redisConfig.yml"
	rabbitConfigFile   = "/values/rabbitConfig.yml"
	consumerConfigFile = "/values/consumerConfig.yml"
)

func main() {
	var redisUrl string
	readConfig(redisFile, &redisUrl)
	common.InitRedisPool(redisUrl)

	var rabbitConfig common.RabbitConfig
	readConfig(rabbitConfigFile, &rabbitConfig)

	var consumerConfigs common.ConsumerConfigs
	readConfig(consumerConfigFile, &consumerConfigs)
	fmt.Println(consumerConfigs)

	setupConsumer(rabbitConfig, consumerConfigs.Configs)

	select {}
}

func readConfig(fileName string, out interface{}) {
	pwd, err := os.Getwd()
	readFile, err := ioutil.ReadFile(pwd + fileName)
	if err != nil {
		log.Fatalf("readfile %s: %w", readFile, err)
	}

	err = yaml.Unmarshal(readFile, out)
	if err != nil {
		log.Fatalf("parse error %s: %w", readFile, err)
	}
}

func setupConsumer(rabbitConfig common.RabbitConfig, consumerConfigs []common.ConsumerConfig) {
	rabbit := common.NewRabbit(rabbitConfig)

	var err error
	if err = rabbit.ConnectRabbit(); err != nil {
		log.Fatalf("unable to connect to rabbit: %w", err)
	}

	for _, consumerConfig := range consumerConfigs {
		consumerConfig.Reconnect.MaxAttempt = 60
		consumerConfig.Reconnect.Interval = 1 * time.Second
		consumer := common.NewConsumer(consumerConfig, rabbit)

		if err = consumer.ConsumerStart(); err != nil {
			log.Fatalf("unable to start consumer: %w", err)
		}
	}
}
