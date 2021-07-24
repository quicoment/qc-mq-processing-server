package common

import (
	"encoding/json"
	"github.com/quicoment/qc-mq-processing-server/domain"
	"log"
	"time"
)

func INSERT_TEST() {
	var c = domain.Comment{1, time.Now(), time.Now(), "content", "password", 0}
	data, _ := json.Marshal(c)
	err := INSERT(1, data)
	if err != nil {
		log.Fatal(err.Error())
	}
}
