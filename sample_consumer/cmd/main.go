package main

import (
	"fmt"
	"net/http"
)

func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Hello Consumer Sample")
	fmt.Fprintf(writer, "Hello Consumer Sample")
}

func health_handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Hello health")
	fmt.Fprintf(writer, "health")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/healthz", health_handler)
	http.ListenAndServe(":8080", nil)
}
