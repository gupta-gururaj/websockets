package main

import (
	"fmt"
	"log"
	"net/http"

	nats "github.com/nats-io/nats.go"

	"ws/internal/handlers"
)

func main() {
	mux := routes()
	log.Println("Starting channel listener")
	nc, _ := nats.Connect(nats.DefaultURL)
	// Creates JetStreamContext
	js, err := nc.JetStream()
	if err != nil {
		fmt.Println("Err::>>>>", err)
	}
	// Creates stream
	var app = &handlers.App{
		Nc: nc,
		Js: js,
	}
	a := []string{"ws_sub1"}
	err = createStream(js, "ws_stream1", a)
	if err != nil {
		fmt.Println("Error crating stream", err)
	}
	go handlers.ListenToWsChannel(app)

	log.Println("Staring web server on port 8083")

	err = http.ListenAndServe(":8083", mux)
	if err != nil {
		fmt.Println("err in starting server", err)
	}
}

func createStream(js nats.JetStreamContext, streamName string, subjects []string) error {
	// Check if the ORDERS stream already exists; if not, create it.
	fmt.Println(">>>>>>>>>", streamName)
	stream, err := js.StreamInfo(streamName)
	if err != nil {
		log.Println(err)
	}
	if stream == nil {
		log.Printf("creating stream %q", streamName, "and subjects my_subject1")
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: subjects,
		})
		if err != nil {
			return err
		}
	} else {
		fmt.Println(">>>>>>>>>??", streamName)
	}
	return nil
}
