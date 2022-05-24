package main

import (
	"io"
	"log"
	"os"
	"strings"

	health_check "github.com/nautible/nautible-kong-serverless/plugin/pkg/health_check"
	pubsub "github.com/nautible/nautible-kong-serverless/plugin/pkg/pubsub"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

type Config struct {
	Backend []struct {
		Target string
		Health string
		Pubsub string
		Topic  string
	}
	Count    int
	Interval int
}

func main() {
	loggingSettings("/tmp/log.txt")
	server.StartServer(func() interface{} {
		return &Config{}
	}, "0.1", 1)
}

func loggingSettings(filename string) {
	logfile, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetOutput(multiLogFile)
}

func (conf Config) Access(kong *pdk.PDK) {
	log.Println("start Access")
	check, err := kong.Request.GetQueryArg("check")
	if err != nil {
		log.Println(err.Error())
		return
	}
	if check == "none" {
		return
	}
	path, err := kong.Request.GetPath()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, backend := range conf.Backend {
		if strings.HasPrefix(path, backend.Target) {
			pubsub.PublishQueue(backend.Pubsub, backend.Topic)

			log.Println("start health_check")
			health_check.Execute(backend.Target, backend.Health, conf.Count, conf.Interval)

		}
	}

	// subscribeQueue(channel, conf)

	log.Println("end Access")

}
