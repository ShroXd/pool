package crawler

import (
	"github.com/gocolly/colly"
	"time"
)

type QuickProxy struct {
	baseURL  string
	domain   string
	cacheDir string
	limit    *colly.LimitRule
}

func (q *QuickProxy) getBaseURL() string {
	return q.baseURL
}

func (q *QuickProxy) getDomain() string {
	return q.domain
}

func (q *QuickProxy) getCacheDir() string {
	return q.cacheDir
}

func (q *QuickProxy) getLimit() *colly.LimitRule {
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
