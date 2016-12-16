package main

import (
	"log"
	"time"

	"github.com/nats-io/go-nats"
)

// 010_OMIT
func main() {
	log.Println("Publisher starting...")
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()
	for {
		nc.Publish("AppA.Stateless.TimeSec", []byte(time.Now().String()))
		time.Sleep(time.Second)
	}
}

// 020_OMIT
