package crawler

import (
	"context"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"log"
	"pool/pkg/db"
	"pool/pkg/model"
)


var ctx = context.Background()
var Colly *colly.Collector

func InitColly() {
	Colly = colly.NewCollector(
		colly.AllowedDomains("free.kuaidaili.com", "kuaidaili.com"),
		colly.CacheDir("./cache"),
		// TODO: disable on prod
		colly.Debugger(&debug.LogDebugger{}),
	)

	Colly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	Colly.OnHTML("tr", func(e *colly.HTMLElement) {

		log.Println("Address ", e.ChildText("td:nth-child(1)"))
		log.Println("Port ", e.ChildText("td:nth-child(2)"))
		log.Println("Anonymous ", e.ChildText("td:nth-child(3)"))
		log.Println("Type ", e.ChildText("td:nth-child(4)"))
		log.Println("Location ", e.ChildText("td:nth-child(5)"))
		log.Println("------------------------")

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
		}

		// TODO: store data from queue async

		err := db.RdbKuaidaili.Set(ctx, addr, agency, 0).Err()
		if err != nil {
			//panic(err)
			log.Println("ERROR!!!!!!1")
			log.Println(err)
		}
	})
}
