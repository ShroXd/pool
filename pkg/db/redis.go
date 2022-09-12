package db

import (
	"github.com/go-redis/redis/v9"
)


var RdbContext *redis.Client

// https://free.kuaidaili.com/free/intr/
var RdbKuaidaili *redis.Client

func InitRedis() {
	// TODO: make the db to be variable

	RdbContext = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:55000",
		Username: "default",
		Password: "redispw",
		DB:       0,
	})

	RdbKuaidaili = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:55000",
		Username: "default",
		Password: "redispw",
		DB:       1,
	})
}
