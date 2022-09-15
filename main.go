package main

import (
	"context"
	"pool/pkg/crawler"
	"pool/pkg/db"
	"pool/pkg/pubsub"
	"time"
)

var ctx = context.Background()
var quit chan int

func main() {
	quit = make(chan int)
	defer func() { <-quit }()

	initDeps()

	p := pubsub.NewPublisher(10*time.Second, 100)
	all := p.Subscribe()

	// TODO: figure out how to run different crawlers concurrently
	//crawler.Run(crawler.NewCloudProxy(), p)
	crawler.Run(crawler.NewQuickProxy(), p)
	go db.WriteData(ctx, all, db.StoreFnBuilder(db.RdbProxy), quit)
}

func initDeps() {
	db.InitRedis()
}
