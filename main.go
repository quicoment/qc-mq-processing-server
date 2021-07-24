package main

import (
	"flag"
	"github.com/quicoment/qc-mq-processing-server/common"
)

var (
	redisServer = flag.String("127.0.0.1", ":6379", "redis-connect-host")
)

func main() {
	common.InitRedisPool(*redisServer)
}
