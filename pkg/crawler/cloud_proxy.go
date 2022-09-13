package crawler

import (
	"context"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"pool/pkg/db"
	"pool/pkg/model"
	"strconv"
	"strings"
	"time"
)

var ctx = context.Background()

type CloudProxy struct {
	baseURL  string
	domain   string
	cacheDir string
	limit    *colly.LimitRule
}

func (c CloudProxy) getBaseURL() string {
	return c.baseURL
}

func (c CloudProxy) getDomain() string {
	return c.domain
}

func (c CloudProxy) getCacheDir() string {
	return c.cacheDir
}

func (c CloudProxy) getLimit() *colly.LimitRule {
	return c.limit
}

func (c CloudProxy) New() CloudProxy {
	cloud := CloudProxy{}

	cloud.baseURL = "http://www.ip3366.net/"
	// TODO: support multi domains
	cloud.domain = "www.ip3366.net"
	cloud.cacheDir = "./cache"
	cloud.limit = &colly.LimitRule{
		DomainGlob:  "*ip3366.*",
		Parallelism: 1,
		Delay:       20 * time.Second,
	}

	return cloud
}

func (c CloudProxy) UrlParser(q *queue.Queue) (string, colly.HTMLCallback) {
	selector := "div[id=listnav]"

	fn := func(e *colly.HTMLElement) {
		// TODO: page url generator
		total := strings.Split(e.ChildText("strong"), "/")[1]

		max, err := strconv.Atoi(total)
		if err != nil {
			log.Println(err)
		}

		for i := 1; i < max; i++ {
			err := q.AddURL(c.baseURL + "?stype=1&page=" + strconv.Itoa(i))
			if err != nil {
				log.Println(err)
			}
		}
	}

	return selector, fn
}

func (c CloudProxy) IpParser() (string, colly.HTMLCallback) {
	selector := "tr"

	fn := func(e *colly.HTMLElement) {
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
	}

	return selector, fn
}
