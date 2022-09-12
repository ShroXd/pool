package main

import (
	"container/list"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"log"
)

type Agency struct {
	Address   string
	Port      string
	Anonymous string
	Type      string
	Location  string
}

func main() {
	c := colly.NewCollector(
		//colly.AllowedDomains("freeproxylists.net", "www.freeproxylists.net"),
		colly.CacheDir("./cache"),
		colly.Debugger(&debug.LogDebugger{}),
	)

	agencies := list.New()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// Find and visit all links
	c.OnHTML("tr", func(e *colly.HTMLElement) {

		log.Println("Address ", e.ChildText("td:nth-child(1)"))
		log.Println("Port ", e.ChildText("td:nth-child(2)"))
		log.Println("Anonymous ", e.ChildText("td:nth-child(3)"))
		log.Println("Type ", e.ChildText("td:nth-child(4)"))
		log.Println("Location ", e.ChildText("td:nth-child(5)"))
		log.Println("------------------------")

		agency := Agency{
			Address: e.ChildText("td:nth-child(1)"),
			Port: e.ChildText("td:nth-child(2)"),
			Anonymous: e.ChildText("td:nth-child(3)"),
			Type: e.ChildText("td:nth-child(4)"),
			Location: e.ChildText("td:nth-child(5)"),
		}

		agencies.PushBack(agency)
	})

	//c.Visit("https://www.freeproxylists.net/zh/")
	c.Visit("https://free.kuaidaili.com/free/intr/")
}
