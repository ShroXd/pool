package crawler

import (
	"context"
	"github.com/gocolly/colly"
	"log"
	"pool/pkg/db"
	"pool/pkg/model"
	"time"
)

var ctx = context.Background()

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

		if err := db.RdbKuaidaili.Set(ctx, addr, agency, 0).Err(); err != nil {
			log.Println("Error during writing IPs: ", err)
		}
	})
}
