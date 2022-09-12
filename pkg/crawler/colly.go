package crawler

import (
	"context"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)


var ctx = context.Background()
var Colly *colly.Collector

func InitColly() {
	Colly = colly.NewCollector(
		colly.AllowedDomains("free.kuaidaili.com", "kuaidaili.com"),
		colly.CacheDir("./cache"),
		// TODO: disable on prod
		colly.Debugger(&debug.LogDebugger{}),
	)

	Colly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	//RegisterIP(Colly)
	RegisterPage(Colly)
}
