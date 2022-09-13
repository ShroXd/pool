package db

import (
	"context"
	"github.com/go-redis/redis/v9"
	"log"
	"pool/pkg/model"
	"pool/pkg/pubsub"
	"time"
)

type StoreFn = func(ctx context.Context, key string, value interface{}, expiration time.Duration) error

func StoreFnBuilder(db *redis.Client) StoreFn {
	return func(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
		if err := db.Set(ctx, key, value, expiration).Err(); err != nil {
			return err
		}

		return nil
	}
}

func WriteData(ctx context.Context, sub pubsub.Subscriber, fn StoreFn, quit chan int) {
	for data := range sub {
		err := fn(ctx, data.(model.Agency).Address, data, 0)
		if err != nil {
			log.Println("Error during writing IPs: ", err)
		}
	}

	quit <- 0
}
