package main

import (
	"github.com/garyburd/redigo/redis"
)

var pool = &redis.Pool{
	MaxIdle:   1,
	MaxActive: 510, // max number of connections
	Dial: func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", ":6379")
		if err != nil {
			panic(err.Error())
		}
		return c, err
	},
}
