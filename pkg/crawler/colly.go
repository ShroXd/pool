package crawler

import (
	"context"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/queue"
	"log"
	"pool/pkg/db"
	"pool/pkg/model"
	"strconv"
	"time"
)

var Colly *colly.Collector
var ctx = context.Background()

const baseURL = "https://free.kuaidaili.com/free/intr/"

func Run() {
	initColly()

	q, _ := queue.New(10, &queue.InMemoryQueueStorage{MaxSize: 10000})
	initUrls(q)

	err := q.Run(Colly)
	if err != nil {
		log.Println(err)
	}
}

func initColly() {
	Colly = colly.NewCollector(
		//colly.AllowedDomains("free.kuaidaili.com", "kuaidaili.com"),
		colly.CacheDir("./cache"),
		// TODO: disable on prod
		colly.Debugger(&debug.LogDebugger{}),
	)

	//rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:4781")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//Colly.SetProxyFunc(rp)

	extensions.RandomUserAgent(Colly)

	Colly.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	Colly.Limit(&colly.LimitRule{
		DomainGlob: "*kuaidaili.*",
		Parallelism: 4,
		RandomDelay: 1 * time.Second,
	})

	registerIP(Colly)
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
			err := q.AddURL(baseURL + strconv.Itoa(i) + "/")
			if err != nil {
				log.Println(err)
			}
		}
	})

	err := urlsCollector.Visit(baseURL)
	if err != nil {
		log.Println(err)
	}
}

func registerIP(c *colly.Collector) {
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		addr := e.ChildText("td:nth-child(1)")
		if addr == "" {
			return
		}

		agency := model.Agency{
			Address:   addr,
			Port:      e.ChildText("td:nth-child(2)"),
			Anonymous: e.ChildText("td:nth-child(3)"),
			Type:      e.ChildText("td:nth-child(4)"),
			Location:  e.ChildText("td:nth-child(5)"),
			Timestamp: time.Now(),
		}

		// TODO: store data from queue async

		if err := db.RdbKuaidaili.Set(ctx, addr, agency, 0).Err(); err != nil {
			log.Println("Error during writing IPs: ", err)
		}
	})
}
