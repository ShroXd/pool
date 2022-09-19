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
	// TODO: queue support different data source and use the related processor

	quit = make(chan int)
	defer func() { <-quit }()

	initDeps()

	p := pubsub.NewPublisher(10*time.Second, 100)
	all := p.Subscribe()

	go db.WriteData(ctx, all, db.StoreFnBuilder(db.RdbProxy), quit)

	// TODO: send summary e-mail
	//crawler.Run(crawler.NewFreeProxy(), p)
	//crawler.Run(crawler.NewCloudProxy(), p)
	//crawler.Run(crawler.NewQuickProxy(), p)
	crawler.Run(crawler.NewHideProxy(), p)
	//crawler.FetchProxy(p)
}

func initDeps() {
	db.InitRedis()
}
