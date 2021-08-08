package common

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/quicoment/qc-mq-processing-server/domain"
	"log"
	"os"
)

var (
	redisPool *redis.Pool
)

func InitRedisPool(address string) {
	redisPool = &redis.Pool{
		MaxIdle:   20,
		MaxActive: 100,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", address)
			if err != nil {
				err = errors.Errorf("ERROR: fail init redis: %w", err)
				log.Fatal(err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
}

func createComment(comment domain.Comment) error {
	conn := redisPool.Get()
	defer conn.Close()

	data, _ := json.Marshal(comment)
	conn.Send("MULTI")
	// register list key = post:{postID}:comment
	conn.Send("LPUSH", "post:"+comment.ID+":comment", comment.ID)
	// like sorted set key = post:{postID}:likes
	conn.Send("ZADD", "post:"+comment.ID+":likes", 0, comment.ID)
	// string key = comment:{commentID}:cache
	conn.Send("SET", "comment:"+comment.ID+":cache", data)
	_, err := redis.Values(conn.Do("EXEC"))

	if err != nil {
		return err
	}

	return nil
}

func likeComment(userId string, postId int64, commentId string) error {
	conn := redisPool.Get()
	defer conn.Close()

	// like userId set key = comment:{commentId}
	isMember, err := redis.Int64(conn.Do("SISMEMBER", "comment:"+commentId, userId))
	if isMember == 1 {
		return errors.Errorf("Same Person like comment: %s", commentId)
	}
	if err != nil {
		return err
	}

	conn.Send("MULTI")
	// like userId set key = comment:{commentId}
	conn.Send("SADD", "comment:"+commentId, userId)
	// like sorted set key = post:{postID}:likes
	conn.Send("ZINCRBY", "post:"+commentId+":likes", 1, commentId)
	_, err = redis.Values(conn.Do("EXEC"))

	if err != nil {
		return err
	}

	return nil
}

func updateComment(comment domain.Comment) error {
	conn := redisPool.Get()
	defer conn.Close()

	// data, _ := json.Marshal(comment)
	// TODO: UPSERT
	conn.Send("MULTI")
	// conn.Send("SET", "comment:"+comment.ID+":cache", data)
	_, err := redis.Values(conn.Do("EXEC"))

	if err != nil {
		return err
	}

	return nil
}
