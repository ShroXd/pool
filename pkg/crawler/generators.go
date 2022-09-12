package crawler

import (
	"log"
	"pool/pkg/db"
)

func GenerateURL() string {
	total, err := db.RdbContext.Get(ctx, "page:total").Result()
	if err != nil {
		log.Println("Error: ", err)
	}

	current, err := db.RdbContext.Get(ctx, "page:current").Result()
	if err != nil {
		log.Println("Error: ", err)
	}

	log.Println("Total page: ", total)
	log.Println("Total page: ", current)

	if current < total {
		res := "https://free.kuaidaili.com/free/intr/" + current + "/"

		next, _ := increaseNum(current)
		log.Println("next: ", next)
		if err := db.RdbContext.Set(ctx, "page:current", next, 0); err != nil {
			log.Println("Error: ", err)
		}

		log.Println("res: ", res)
		return res
	} else {
		return ""
	}
}
