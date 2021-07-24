package common

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"log"
)

func connection() redis.Conn {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err.Error())
	}

	pong, err := redis.String(c.Do("PING"))
	if err != nil {
		c.Close()
		log.Fatal(err.Error())
	}
	fmt.Printf("PING Response = %s\n", pong)
	return c
}

func GET(key int64) ([]byte, error) {
	conn := connection()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, errors.Errorf("error getting key %d: %w", key, err)
	}
	return data, err
}

func GETALL(pattern string) ([]byte, error) {
	conn := connection()
	defer conn.Close()

	var keys []int64
	var data []byte
	keys, err := redis.Int64s(conn.Do("KEYS", pattern))

	if err != nil {
		return nil, errors.Errorf("error getting %s: %w", pattern, err)
	}

	for _, key := range keys {
		var d, _ = redis.String(conn.Do("GET", key))
		data = append(data, d...)
	}

	return data, nil
}

func INSERT(key int64, value []byte) error {
	conn := connection()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)

	if err != nil {
		return errors.Errorf("error set key %d: %w", key, err)
	}
	return err
}

func DELETE(key int64) error {
	conn := connection()
	defer conn.Close()

	_, err := conn.Do("DEL", key)

	if err != nil {
		return errors.Errorf("error delete key %d: %w", key, err)
	}
	return err
}

func Incr(counterKey string) (int64, error) {
	conn := connection()
	defer conn.Close()

	key, err := redis.Int64(conn.Do("INCR", counterKey))

	if err != nil {
		return -1, errors.Errorf("error get increment key of %s: %w", counterKey, err)
	}
	return key, err
}

func INSERT_SET(setName string, value string) error {
	conn := connection()
	defer conn.Close()

	finished, err := redis.Int(conn.Do("SADD", setName, value))
	if err != nil {
		return errors.Errorf("%d commands were successful, but not completed: %w", finished, err)
	}
	return err
}

func DELETE_SET(setName string, value string) error {
	conn := connection()
	defer conn.Close()

	finished, err := redis.Int(conn.Do("SREM", setName, value))
	if err != nil {
		return errors.Errorf("%d commands were successful, but not completed: %w", finished, err)
	}
	return err
}

func GETALL_SET(setName string) ([]byte, error) {
	conn := connection()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("SMEMBERS", setName))

	if err != nil {
		return nil, errors.Errorf("error get key %s: %w", setName, err)
	}
	return data, err
}
