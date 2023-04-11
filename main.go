package main

import (
	"MMDB/handler"
	"github.com/rs/cors"
	"log"
	"net/http"
	"time"
)

func main() {
	var serveMux = http.NewServeMux()
	serveMux.HandleFunc("/initDatabase", handler.InitDatabase)
	serveMux.HandleFunc("/", handler.HelloWorld)
	handler := cors.Default().Handler(serveMux)
	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       3 * time.Second,
		Handler:           handler,
	}
	log.Fatal(server.ListenAndServe())
}
