package crawler

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"pool/pkg/model"
	"pool/pkg/pubsub"
	"strconv"
	"time"
)

type QuickProxy struct {
	baseURL  string
	domain   string
	cacheDir string
	limit    *colly.LimitRule
}

func (q QuickProxy) getBaseURL() string {
	return q.baseURL
}

func (q QuickProxy) getDomain() string {
	return q.domain
}

func (q QuickProxy) getCacheDir() string {
	return q.cacheDir
}

func (q QuickProxy) getLimit() *colly.LimitRule {
	return q.limit
}

func NewQuickProxy() QuickProxy {
	q := QuickProxy{}

	q.baseURL = "https://free.kuaidaili.com/free/inha/"
	q.domain = "free.kuaidaili.com"
	q.cacheDir = "./cache"
	q.limit = &colly.LimitRule{
		DomainGlob:  "*.kuaidaili.*",
		Parallelism: 1,
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	}

	return q
}

func (q QuickProxy) UrlParser(queue *queue.Queue) (string, colly.HTMLCallback) {
	// TODO: also collect proxy form https://free.kuaidaili.com/free/intr/
	selector := "div[id=listnav]"

	fn := func(e *colly.HTMLElement) {
		total := e.ChildText("li:nth-last-child(2)")

		max, err := strconv.Atoi(total)
		if err != nil {
			log.Println(err)
		}

		log.Println("Max: ", max)
		for i := 1; i < 2; i++ {
			log.Println("Add url to queue: ", q.baseURL + strconv.Itoa(i))
			err := queue.AddURL(q.baseURL + strconv.Itoa(i))
			if err != nil {
				log.Println(err)
			}
		}
	}

	return selector, fn
}

func (QuickProxy) IpParser(p *pubsub.Publisher) (string, colly.HTMLCallback) {
	selector := "tr"

	fn := func(e *colly.HTMLElement) {
		addr := e.ChildText("td:nth-child(1)")
		if addr == "" {
			return
		}

		agency := model.Agency{
			Address: addr,
			Port: e.ChildText("td:nth-child(2)"),
			Anonymous: e.ChildText("td:nth-child(3)"),
			Type:      e.ChildText("td:nth-child(4)"),
			Location:  e.ChildText("td:nth-child(5)"),
			Timestamp: time.Now(),
		}

		p.Publish(agency)
	}

	return selector, fn
}
