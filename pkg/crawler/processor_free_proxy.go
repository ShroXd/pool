package crawler

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"pool/pkg/model"
	"pool/pkg/pubsub"
	"time"
)

type FreeProxy struct {
	baseURL  string
	domain   string
	cacheDir string
	limit    *colly.LimitRule
}

func (f FreeProxy) getBaseURL() string {
	return f.baseURL
}

func (f FreeProxy) getDomain() string {
	return f.domain
}

func (f FreeProxy) getCacheDir() string {
	return f.cacheDir
}

func (f FreeProxy) getLimit() *colly.LimitRule {
	return f.limit
}

func NewFreeProxy() FreeProxy {
	proxy := FreeProxy{}

	proxy.baseURL = "https://free-proxy-list.net/"
	proxy.domain = "free-proxy-list.net"
	proxy.cacheDir = "./cache"
	proxy.limit = &colly.LimitRule{
		DomainGlob:  "*free-proxy-list.*",
		Parallelism: 1,
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	}

	return proxy
}

func (f FreeProxy) UrlParser(q *queue.Queue) (string, colly.HTMLCallback) {
	selector := "html"

	fn := func(e *colly.HTMLElement) {
		// TODO: refactor the colly function to receive worker func
		err := q.AddURL(f.baseURL)
		if err != nil {
			log.Println(err)
		}
	}

	return selector, fn
}

func (f FreeProxy) IpParser(p *pubsub.Publisher) (string, colly.HTMLCallback) {
	selector := "tr"

	fn := func(e *colly.HTMLElement) {
		addr := e.ChildText("td:nth-child(1)")
		if !isValidIp(addr) {
			return
		}

		// TODO: use protocol instead of type
		agency := model.Agency{
			Address:   addr,
			Port:      e.ChildText("td:nth-child(2)"),
			Anonymous: e.ChildText("td:nth-child(5)"),
			Type:      proxyType(e.ChildText("td:nth-child(7)")),
			Location:  e.ChildText("td:nth-child(4)"),
			Timestamp: time.Now(),
		}

		p.Publish(agency)
	}

	return selector, fn
}
