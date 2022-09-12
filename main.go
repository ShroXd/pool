package main

import (
	"pool/pkg/crawler"
	"pool/pkg/db"
)

func main() {
	initDeps()

	crawler.Run()
}

func initDeps() {
	db.InitRedis()
}
