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
var quit chan int

func main() {
	quit = make(chan int)

	initDeps()

	p := pubsub.NewPublisher(10*time.Second, 100)
	all := p.Subscribe()

	crawler.Run(crawler.CloudProxy{}.New(), p)

	go func() {
		for agency := range all {
			if err := db.RdbProxy.Set(ctx, agency.(model.Agency).Address, agency, 0).Err(); err != nil {
				log.Println("Error during writing IPs: ", err)
			}
		}

		quit <- 0
	}()

	<-quit
}

func initDeps() {
	db.InitRedis()
}
