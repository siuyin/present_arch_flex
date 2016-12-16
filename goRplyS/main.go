package main

import (
	"log"
	"time"

	"github.com/nats-io/go-nats"
)

// 010_OMIT
func main() {
	log.Println("Slow Replier (Server) starting...")
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	nc.Subscribe("AppA.Stateless.HelloAdder", func(m *nats.Msg) {
		time.Sleep(5 * time.Millisecond)
		nc.Publish(m.Reply, []byte("Hello S: "+string(m.Data))) // HL
	})

	select {} // wait forever
}

// 020_OMIT
