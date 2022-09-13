package main

import (
	"fmt"
	"pool/pkg/pubsub"
	"strings"
	"time"
)

func main() {
	//initDeps()
	//
	//crawler.Run(crawler.CloudProxy{}.New())

	p := pubsub.NewPublisher(10*time.Second, 100)
	p.Close()

	all := p.Subscribe()
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}

		return false
	})

	p.Publish("Hello, world!")
	p.Publish("Hello, golang!")

	go func() {
		for msg := range all {
			fmt.Println("all:", msg)
		}
	}()

	go func() {
		for msg := range golang {
			fmt.Println("golang:", msg)
		}
	}()

	time.Sleep(10 * time.Second)
}

//func initDeps() {
//	db.InitRedis()
//}
