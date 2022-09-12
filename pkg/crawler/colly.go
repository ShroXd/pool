package crawler

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"strconv"
)

var Colly *colly.Collector

func InitColly() {
	baseURL := "https://free.kuaidaili.com/free/intr/"

	Colly = colly.NewCollector(
		colly.AllowedDomains("free.kuaidaili.com", "kuaidaili.com"),
		colly.CacheDir("./cache"),
		// TODO: disable on prod
		//colly.Debugger(&debug.LogDebugger{}),
	)

	q, _ := queue.New(10, &queue.InMemoryQueueStorage{MaxSize: 10000})

	Colly.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	RegisterIP(Colly)
	//RegisterPage(Colly)

	for i := 1; i < 10; i++ {
		q.AddURL(baseURL + strconv.Itoa(i) + "/")
	}

	q.Run(Colly)
}
