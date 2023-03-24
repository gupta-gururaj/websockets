package main

import (
	"fmt"
	"log"
	"net/http"
	"ws/internal/handlers"
)

func main() {
	mux := routes()

	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	log.Println("Staring web server on port 8083")

	err := http.ListenAndServe(":8083", mux)
	if err != nil {
		fmt.Println("err in starting server", err)
	}
}
