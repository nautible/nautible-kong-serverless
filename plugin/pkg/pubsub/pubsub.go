package pubsub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Ping struct {
	Id int `json: "id"`
}

// Daprを経由してPub/Subを行う
func PublishQueue(pubsub string, topic string) {
	log.Println("start PublishQueue")

	// キューに入れるデータを生成（トリガーなので何でもよい）
	ping := new(Ping)
	ping.Id = 1
	ping_json, _ := json.Marshal(ping)

	url_target := fmt.Sprintf("http://localhost:3500/v1.0/publish/%s/%s", pubsub, topic)
	log.Println(url_target)
	log.Println("start Post")
	res, err := http.Post(url_target, "application/json", bytes.NewBuffer(ping_json))
	if err != nil {
		log.Println("Request error:", err)
		return
	}
	defer res.Body.Close()
	log.Println("end Post")
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Request error:", err)
		return
	}

	str_json := string(body)
	log.Println(str_json)
}
