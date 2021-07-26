package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/gin-contrib/cors"
	"github.com/quicoment/qc-mq-processing-server/api"
	"github.com/quicoment/qc-mq-processing-server/common"
	"time"
)

var (
	redisServer = flag.String("127.0.0.1", ":6379", "redis-connect-host")
)

func main() {
	common.InitRedisPool(*redisServer)

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

	r.GET("/queues", api.CreateQueue)

	return r
}
