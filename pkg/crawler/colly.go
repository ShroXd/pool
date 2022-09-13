package crawler

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/queue"
	"log"
)

var Colly *colly.Collector

func Run(cp ProxyWebsite) {
	initColly(cp)

	q, _ := queue.New(10, &queue.InMemoryQueueStorage{MaxSize: 10000})
	initUrls(cp, q)

	err := q.Run(Colly)
	if err != nil {
		log.Println(err)
	}
}

func initColly(cp ProxyWebsite) {
	Colly = colly.NewCollector(
		colly.AllowedDomains(cp.getDomain()),
		colly.CacheDir(cp.getCacheDir()),
		// TODO: disable on prod
		colly.Debugger(&debug.LogDebugger{}),
	)

	// TODO: add proxy for each request

	extensions.RandomUserAgent(Colly)

	Colly.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})
	Colly.Limit(cp.getLimit())
	Colly.OnHTML(cp.IpParser())
}

func initUrls(cp ProxyWebsite, q *queue.Queue) {
	urlsCollector := Colly.Clone()
	urlsCollector.OnHTML(cp.UrlParser(q))

	err := urlsCollector.Visit(cp.getBaseURL())
	if err != nil {
		log.Println(err)
	}
}
