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
	"strings"
	"time"
)

var Colly *colly.Collector
var ctx = context.Background()

const baseURL = "http://www.ip3366.net/"

// TODO: encapsulate the logic for single website

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
		colly.AllowedDomains("www.ip3366.net"),
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
		DomainGlob:  "*ip3366.*",
		Parallelism: 1,
		Delay:       20 * time.Second,
	})

	registerIP(Colly)
}

func initUrls(q *queue.Queue) {
	urlsCollector := Colly.Clone()

	urlsCollector.OnHTML("div[id=listnav]", func(e *colly.HTMLElement) {
		// TODO: page url generator
		total := strings.Split(e.ChildText("strong"), "/")[1]

		max, err := strconv.Atoi(total)
		if err != nil {
			log.Println(err)
		}

		for i := 1; i < max; i++ {
			err := q.AddURL(baseURL + "?stype=1&page=" + strconv.Itoa(i))
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
			Anonymous: convertChinese(e.ChildText("td:nth-child(3)")),
			Type:      e.ChildText("td:nth-child(4)"),
			Location:  convertChinese(e.ChildText("td:nth-child(6)")),
			Timestamp: time.Now(),
		}

		// TODO: store data from queue async

		if err := db.RdbProxy.Set(ctx, addr, agency, 0).Err(); err != nil {
			log.Println("Error during writing IPs: ", err)
		}
	})
}
