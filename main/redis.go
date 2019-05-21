package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func setKey(resource res) interface{} {
	c, err := redis.Dial("tcp", "localhost: 6379")
	if err != nil {
		fmt.Println("Connection to redis failed:", err)
		return nil
	}
	defer c.Close()
	return nil
}
