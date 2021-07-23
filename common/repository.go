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

func Get(key int64) ([]byte, error) {
	conn := connection()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, err
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

func POST(value []byte) error {
	conn := connection()
	defer conn.Close()

	var key, err = Incr("COMMENT")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

func Exists(key int64) (bool, error) {
	conn := connection()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

func Delete(key int64) error {
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
