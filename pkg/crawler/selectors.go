package crawler

import (
	"github.com/gocolly/colly"
	"log"
	"pool/pkg/db"
	"pool/pkg/model"
	"time"
)

func RegisterIP(c *colly.Collector) {
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

		err := db.RdbKuaidaili.Set(ctx, addr, agency, 0).Err()
		if err != nil {
			//panic(err)
			log.Println("ERROR!!!!!!1")
			log.Println(err)
		}
	})
}

func RegisterPage(c *colly.Collector) {
	c.OnHTML("div[id=listnav]", func(e *colly.HTMLElement) {
		// TODO: Store it in the redis
		// TODO: page url generator
		lastPage := e.ChildText("li:nth-last-child(2)")
		log.Println("lastPage: ", lastPage)
	})
}
