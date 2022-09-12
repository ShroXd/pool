package main

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
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
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

var ctx = context.Background()

func (t Agency) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t Agency) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	return nil
}

func main() {
	c := colly.NewCollector(
		//colly.AllowedDomains("freeproxylists.net", "www.freeproxylists.net"),
		colly.CacheDir("./cache"),
		colly.Debugger(&debug.LogDebugger{}),
	)

	//agencies := list.New()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:55000",
		Username: "default",
		Password: "redispw",
		DB:       0,
	})

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

		addr := e.ChildText("td:nth-child(1)")
		if addr == "" {
			return
		}

		agency := Agency{
			Address:   addr,
			Port:      e.ChildText("td:nth-child(2)"),
			Anonymous: e.ChildText("td:nth-child(3)"),
			Type:      e.ChildText("td:nth-child(4)"),
			Location:  e.ChildText("td:nth-child(5)"),
		}

		// store data from queue async
		//agencies.PushBack(agency)

		err := rdb.Set(ctx, addr, agency, 0).Err()
		if err != nil {
			//panic(err)
			log.Println("ERROR!!!!!!1")
			log.Println(err)
		}
	})

	//c.Visit("https://www.freeproxylists.net/zh/")
	c.Visit("https://free.kuaidaili.com/free/intr/")
}
