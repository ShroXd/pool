package crawler

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"pool/pkg/pubsub"
)

type ProxyWebsite interface {
	getBaseURL() string
	getDomain() string
	getCacheDir() string
	getLimit() *colly.LimitRule

	UrlParser(q *queue.Queue) (string, colly.HTMLCallback)
	IpParser(p *pubsub.Publisher) (string, colly.HTMLCallback)
}
