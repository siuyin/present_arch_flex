package main

import (
	"fmt"
	"log"

	"github.com/nats-io/go-nats"
)

// 010_OMIT
func main() {
	log.Println("Subscriber starting...")
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()
	// Try: Time*, Statel> and AppA.* and AppA.> // HL
	nc.Subscribe("AppA.Stateless.TimeSec", func(m *nats.Msg) {

		fmt.Printf("%s: %s\n", m.Subject, string(m.Data))
	})
	select {} // wait forever
}

// 020_OMIT
