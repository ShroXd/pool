package main

import (
	"context"
	"log"
	"pool/pkg/crawler"
	"pool/pkg/db"
	"pool/pkg/model"
	"pool/pkg/pubsub"
	"time"
)

var ctx = context.Background()

func main() {
	initDeps()

	p := pubsub.NewPublisher(10*time.Second, 100)
	// TODO: How to control the close for multi crawler
	defer p.Close()

	all := p.Subscribe()
	crawler.Run(crawler.CloudProxy{}.New(), p)

	go func() {
		for agency := range all {
			if err := db.RdbProxy.Set(ctx, agency.(model.Agency).Address, agency, 0).Err(); err != nil {
				log.Println("Error during writing IPs: ", err)
			}
		}
	}()

	// TODO: close the channel correctly instead of exit main function directly
}

func initDeps() {
	db.InitRedis()
}
