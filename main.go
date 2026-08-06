package main

import (
	"log"
	"net/http"
)

func main() {
	//fmt.Println("Hello, World!")
	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
