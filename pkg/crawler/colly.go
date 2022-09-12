package crawler

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"pool/pkg/db"
	"strconv"
)

var Colly *colly.Collector
var baseURL = "https://free.kuaidaili.com/free/intr/"

func InitColly() {

	Colly = colly.NewCollector(
		colly.AllowedDomains("free.kuaidaili.com", "kuaidaili.com"),
		colly.CacheDir("./cache"),
		// TODO: disable on prod
		//colly.Debugger(&debug.LogDebugger{}),
	)

	q, _ := queue.New(10, &queue.InMemoryQueueStorage{MaxSize: 10000})
	initUrls(q)

	Colly.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	RegisterIP(Colly)

	q.Run(Colly)
}

func initUrls(q *queue.Queue) {
	urlsCollector := Colly.Clone()

	urlsCollector.OnHTML("div[id=listnav]", func(e *colly.HTMLElement) {
		// TODO: page url generator
		total := e.ChildText("li:nth-last-child(2)")

		if err := db.RdbContext.Set(ctx, "page:total", total, 0).Err(); err != nil {
			log.Println("Error during writing page:total ", err)
		}

		if err := db.RdbContext.Set(ctx, "page:current", 1, 0).Err(); err != nil {
			log.Println("Error during writing page:current ", err)
		}

		max, err := strconv.Atoi(total)
		if err != nil {
			log.Println(err)
		}

		for i := 1; i < max; i++ {
			q.AddURL(baseURL + strconv.Itoa(i) + "/")
		}
	})

	urlsCollector.Visit(baseURL)
}
