package main

import (
	"pool/pkg/crawler"
	"pool/pkg/db"
)

func main() {
	initDeps()

	crawler.Run(crawler.CloudProxy{}.New())
}

func initDeps() {
	db.InitRedis()
}
