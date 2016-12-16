package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/go-nats"
)

// 010_OMIT
type pipeCfg struct {
	NumJobs, NumWorkers int // HL
}
type workerRes struct {
	Result string
}

// 020_OMIT
// 030_OMIT
var wg = &sync.WaitGroup{}

func main() {
	log.Println("Pipeline Dispatcher Fast (RRF) starting...")
	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, "json") // or gob
	defer c.Close()
	// 040_OMIT
	//050_OMIT
	numWorkers := 1 // Try 2,5,10,100 // HL
	//060_OMIT
	//070_OMIT
	rrDispatcher() // HL
	wg.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		worker(i)
	}
	//080_OMIT
	//090_OMIT
	numJobs := 10
	//092_OMIT
	for ready := false; !ready; {
		c.Publish("AppA.Stateless.PipeA.Config", pipeCfg{numJobs, numWorkers})
		rdy := "Ready?"
		time.Sleep(300 * time.Microsecond)                              // let the data get there
		if err := c.Request("AppA.Stateful.PipeA.RRDispatcher.IsReady", // HL
			&rdy, &rdy, 10*time.Millisecond); err != nil {
			fmt.Println(time.Now(), rdy)
		} else {
			fmt.Println(time.Now(), rdy)
			ready = true
		}
	}
	//094_OMIT
	for i := 1; i <= numJobs; i++ {
		c.Publish("AppA.Stateful.PipeA.RRDispatcher", &i)
	}
	wg.Wait() // wait for goroutines to finish
	c.Flush()
	//100_OMIT
}

//110_OMIT
func worker(i int) {
	go func() {
		done := make(chan bool)
		goneHome := false
		jobQ := make(chan int, 1000)
		nc, _ := nats.Connect(nats.DefaultURL)
		c, _ := nats.NewEncodedConn(nc, "json")
		defer c.Close()
		c.Subscribe(fmt.Sprintf("AppA.Stateless.PipeA.Worker%d", i), func(j *int) { // HL
			jobQ <- *j
		})
		c.Subscribe("AppA.Stateless.PipeA.StopWork", func(s *string) {
			done <- true
		})
		for {
			select {
			case j := <-jobQ:
				time.Sleep(100 * time.Millisecond)
				msg := fmt.Sprintf("Worker%03d: Job: %02d done.\n", i, j)
				c.Publish("AppA.Stateful.PipeA.Cons", workerRes{msg})
				fmt.Printf("Worker%03d sent job %02d results.\n", i, j)
			case <-done:
				if !goneHome {
					wg.Done()
					goneHome = true
				}
			}
		}
	}()
}

//120_OMIT
//130_OMIT
func rrDispatcher() {
	go func() {
		jobQ := make(chan int, 1000)
		nc, _ := nats.Connect(nats.DefaultURL)
		c, _ := nats.NewEncodedConn(nc, "json")
		defer c.Close()
		pc := &pipeCfg{}
		c.Subscribe("AppA.Stateless.PipeA.Config", func(cfg *pipeCfg) { // HL
			pc = cfg // Critical Dispatcher configuration!! // HL
		})
		c.Subscribe("AppA.Stateful.PipeA.RRDispatcher", func(j *int) { // HL
			jobQ <- *j
		})
		//132_OMIT
		c.Subscribe("AppA.Stateful.PipeA.RRDispatcher.IsReady",
			func(subj, reply string, m *string) {
				if pc.NumWorkers != 0 {
					c.Publish(reply, "Yes") // HL
				}
			})
		//134_OMIT
		//135_OMIT
		for {
			select {
			case j := <-jobQ:
				n := j%pc.NumWorkers + 1
				c.Publish(fmt.Sprintf("AppA.Stateless.PipeA.Worker%d", n), &j) // HL
			}
		}
	}()
}

//140_OMIT
