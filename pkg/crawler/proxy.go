package crawler

import "github.com/gocolly/colly"

type ProxyWebsite interface {
	UrlParser() (string, colly.HTMLCallback)
	IpParser() (string, colly.HTMLCallback)
}
