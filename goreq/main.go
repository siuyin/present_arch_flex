package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/go-nats"
)

// 010_OMIT
func main() {
	log.Println("Requester (Client) starting...")
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	m, err := nc.Request("AppA.Stateless.HelloAdder",
		[]byte("AB CD"), 10*time.Millisecond) // Try Nano seconds // HL
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s\n", m.Subject, string(m.Data))
}

// 020_OMIT
