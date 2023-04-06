package main

import (
	"net/http"
)

func main() {
	// create a new server mux
	mux := http.NewServeMux()
	// register the handler function
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	// create a new http server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	// start the server as a goroutine
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
