package main

import (
	"log"
	"time"

	"github.com/nats-io/go-nats"
)

// 010_OMIT
func main() {
	log.Println("Replier (Server) starting...")
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	nc.Subscribe("AppA.Stateless.HelloAdder", func(m *nats.Msg) {
		time.Sleep(1 * time.Millisecond)
		nc.Publish(m.Reply, []byte("Hello: "+string(m.Data))) // HL
	})

	select {} // wait forever
}

// 020_OMIT
