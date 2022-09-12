package main

import (
	"pool/pkg/crawler"
	"pool/pkg/db"
)

func main() {
	initDeps()

	crawler.Colly.Visit("https://free.kuaidaili.com/free/intr/")
}

func initDeps() {
	db.InitRedis()
	crawler.InitColly()
}
