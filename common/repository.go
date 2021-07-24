package common

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
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
		return data, fmt.Errorf("error getting key %s: %v", key, err)
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
		return nil, fmt.Errorf("error getting %s: %v", pattern, err)
	}

	for _, key := range keys {
		var d, _ = redis.String(conn.Do("GET", key))
		data = append(data, d...)
	}

	return data, nil
}

func UPDATE(key int64, value []byte) error {
	conn := connection()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

func POST(key int64, value []byte) error {
	conn := connection()
	defer conn.Close()

	var _, err = conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

func DELETE(key int64) error {
	conn := connection()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func Incr(counterKey string) (int64, error) {
	conn := connection()
	defer conn.Close()

	key, err := redis.Int64(conn.Do("INCR", counterKey))
	return key, err
}

func POSTSet(setName string, value string) error {
	conn := connection()
	defer conn.Close()

	_, err := conn.Do("SADD", setName, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", setName, v, err)
	}
	return err
}

func DELETESet(setName string, value string) error {
	conn := connection()
	defer conn.Close()

	_, err := conn.Do("SPOP", setName, value)
	return err
}

func GETALLSet(setName string) ([]byte, error) {
	conn := connection()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("SMEMBERS", setName))
	return data, err
}
