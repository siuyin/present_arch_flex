package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/go-nats"
)

// 010_OMIT
type pipeCfg struct {
	NumJobs int
}
type workerRes struct {
	Result string
}

// 020_OMIT
// 030_OMIT
func main() {
	log.Println("Pipeline Results Consolidator starting...")
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, "json") // or gob // HL
	defer c.Close()
	// 040_OMIT
	//050_OMIT
	var (
		pc                 pipeCfg
		results            []string
		startTime, endTime time.Time
	)
	//060_OMIT
	//070_OMIT
	c.Subscribe("AppA.Stateless.PipeA.Config", func(p *pipeCfg) {
		pc = *p
		fmt.Printf("Configuration: %v jobs\n", pc.NumJobs)
		startTime = time.Now()
	})
	//080_OMIT
	//090_OMIT
	n := 0
	c.Subscribe("AppA.Stateful.PipeA.Cons", func(r *workerRes) {
		n++
		results = append(results, r.Result)
		if n == pc.NumJobs {
			endTime = time.Now()
			fmt.Println(results)
			fmt.Printf("Processing time: %v\n\n", endTime.Sub(startTime).Seconds())
			n = 0
			results = []string{}
		}
	})
	select {}

	//100_OMIT
}
