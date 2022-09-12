package main

import (
	"pool/pkg/crawler"
	"pool/pkg/db"
)

func main() {
	initDeps()
}

func initDeps() {
	db.InitRedis()
	crawler.InitColly()
}
