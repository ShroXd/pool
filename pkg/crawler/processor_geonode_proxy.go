package crawler

import (
	"github.com/go-resty/resty/v2"
	"log"
	"pool/pkg/model"
	"pool/pkg/pubsub"
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

	res := &Response{}

	_, err := client.R().
		SetResult(res).
		EnableTrace().
		Get("https://proxylist.geonode.com/api/proxy-list?limit=50&page=1&sort_by=lastChecked&sort_type=desc")
	if err != nil {
		log.Println("Err: ", res)
	}

	for _, client := range res.Data {
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

	p.Close()
}
