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

type HideProxy struct {
	baseURL  string
	domain   string
	cacheDir string
	limit    *colly.LimitRule
}

func (h HideProxy) getBaseURL() string {
	return h.baseURL
}

func (h HideProxy) getDomain() string {
	return h.domain
}

func (h HideProxy) getCacheDir() string {
	return h.cacheDir
}

func (h HideProxy) getLimit() *colly.LimitRule {
	return h.limit
}

func NewHideProxy() HideProxy {
	h := HideProxy{}

	h.baseURL = "https://hidemy.name/en/proxy-list/"
	h.domain = "hidemy.name"
	h.cacheDir = "./cache"
	h.limit = &colly.LimitRule{
		DomainGlob:  "*hidemy.*",
		Parallelism: 1,
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	}

	return h
}

func (h HideProxy) UrlParser(q *queue.Queue) (string, colly.HTMLCallback) {
	selector := "div[class=pagination]"

	fn := func(e *colly.HTMLElement) {
		total := e.ChildText("li:nth-last-child(2)")

		max, err := strconv.Atoi(total)
		if err != nil {
			log.Println(err)
		}

		initialNumber := 0

		log.Println("Max: ", max)
		for i := 1; i < max; i++ {
			url := h.baseURL + "?start=" + strconv.Itoa(initialNumber) + "#list"

			err := q.AddURL(url)
			if err != nil {
				log.Println(err)
			}

			initialNumber += 64
		}
	}

	return selector, fn
}

func (HideProxy) IpParser(p *pubsub.Publisher) (string, colly.HTMLCallback) {
	selector := "tr"

	fn := func(e *colly.HTMLElement) {
		addr := e.ChildText("td:nth-child(1)")
		if !isValidIp(addr) {
			return
		}

		agency := model.Agency{
			Address:   addr,
			Port:      e.ChildText("td:nth-child(2)"),
			Anonymous: e.ChildText("td:nth-last-child(2)"),
			Type:      e.ChildText("td:nth-last-child(3)"),
			Location:  e.ChildText("td:nth-child(3)"),
			Timestamp: time.Now(),
		}

		p.Publish(agency)
	}

	return selector, fn
}
