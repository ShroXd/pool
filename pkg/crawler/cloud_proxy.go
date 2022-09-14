package crawler

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"pool/pkg/model"
	"pool/pkg/pubsub"
	"strconv"
	"strings"
	"time"
)

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
		// TODO: generate it based on baseURL
		DomainGlob:  "*.ip3366.*",
		Parallelism: 1,
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
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

func (c CloudProxy) IpParser(p *pubsub.Publisher) (string, colly.HTMLCallback) {
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

		p.Publish(agency)
	}

	return selector, fn
}
