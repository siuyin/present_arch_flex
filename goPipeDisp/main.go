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
	log.Println("Pipeline Dispatcher starting...")
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, "json") // or gob
	defer c.Close()
	// 040_OMIT
	//050_OMIT
	numWorkers := 1 // Try 2,5,10,100 // HL
	//060_OMIT
	//070_OMIT
	for i := 1; i <= numWorkers; i++ {
		worker(i)
	}
	time.Sleep(time.Second) // wait for workers to be ready
	//080_OMIT
	//090_OMIT
	numJobs := 10
	c.Publish("AppA.Stateless.PipeA.Config", pipeCfg{numJobs})
	for i := 1; i <= numJobs; i++ {
		c.Publish("AppA.Stateless.PipeA.Work", &i)
	}
	time.Sleep(2 * time.Second) // wait for goroutines to finish
	//100_OMIT
}

//110_OMIT
func worker(i int) {
	go func() {
		jobQ := make(chan int, 1000)
		nc, _ := nats.Connect(nats.DefaultURL)
		c, _ := nats.NewEncodedConn(nc, "json")
		defer c.Close()
		c.QueueSubscribe("AppA.Stateless.PipeA.Work", "workerPool", func(j *int) { // HL
			jobQ <- *j
		})
		for {
			select {
			case j := <-jobQ:
				time.Sleep(100 * time.Millisecond)
				msg := fmt.Sprintf("Worker%03d: Job: %02d done.\n", i, j)
				c.Publish("AppA.Stateful.PipeA.Cons", workerRes{msg}) // HL
				fmt.Printf("Worker%03d sent job %02d results.\n", i, j)
			}
		}
	}()
}

//120_OMIT
