package db

import (
	"github.com/go-redis/redis/v9"
)


var RdbContext *redis.Client

// https://free.kuaidaili.com/free/intr/
var RdbProxy *redis.Client

func InitRedis() {
	// TODO: make the db to be variable

	RdbContext = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Username: "default",
		Password: "redispw",
		DB:       0,
	})

	RdbProxy = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Username: "default",
		Password: "redispw",
		DB:       1,
	})
}
