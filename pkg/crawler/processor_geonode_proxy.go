package crawler

import (
	"github.com/go-resty/resty/v2"
	"log"
	"pool/pkg/model"
	"pool/pkg/pubsub"
	"strconv"
	"time"
)

type Proxy struct {
	Ip             string
	Port           string
	AnonymityLevel string
	Protocols      []string
	Country        string
	City           string
}

type Response struct {
	Data []*Proxy
}

func FetchProxy(p *pubsub.Publisher) {
	client := resty.New()
	resp := &Response{}

	page := 1
	for {
		_, err := client.R().
			SetResult(resp).
			EnableTrace().
			Get("https://proxylist.geonode.com/api/proxy-list?limit=400&page=" + strconv.Itoa(page) + "&sort_by=lastChecked&sort_type=desc")
		if err != nil {
			log.Println("Err: ", resp)
		}

		log.Println("Fetch data size: ", len(resp.Data))

		if len(resp.Data) == 0 {
			break
		} else {
			for _, client := range resp.Data {
				agency := model.Agency{
					Address:   client.Ip,
					Port:      client.Port,
					Anonymous: client.AnonymityLevel,
					Type:      client.Protocols[0],
					Location:  client.Country + " " + client.City,
					Timestamp: time.Now(),
				}

				log.Println("agency: ", agency)

				p.Publish(agency)
			}
		}

		page++
	}

	p.Close()
}
