package db

import (
	"context"
	"github.com/go-redis/redis/v9"
)

var Ctx = context.Background()

// https://free.kuaidaili.com/free/intr/
var RdbKuaidaili *redis.Client

func InitRedis()  {
	RdbKuaidaili = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:55000",
		Username: "default",
		Password: "redispw",
		// TODO: make the db to be variable
		DB:       1,
	})
}
