package health_check

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Execute(target string, health string, check_count int, interval int) {

	if health == "" {
		return
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for i := 0; i < check_count; i++ {
		url_target := fmt.Sprintf("http://localhost:8000%s%s?check=none", target, health)
		log.Println("start Health Check : " + url_target)
		resp, err := client.Get(url_target)
		if err != nil {
			log.Printf("%s\n", err)
			log.Println("connection error sleep...")
			time.Sleep(time.Duration(interval) * time.Millisecond)
			continue
		}
		log.Println("status : " + resp.Status)
		defer resp.Body.Close()
		if resp.Status == "200 OK" {
			break
		}
		log.Println("application error sleep...")
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}
