package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	w "workers/pkg/worker"
)

func main() {
	numOfWorkers := uint(5)
	channelBuffer := uint(2056)

	workers := w.NewWorkers(numOfWorkers, channelBuffer)
	workers.Run()

	//create a listener to queued jobs
	go func() {
		for err := range workers.PollJob() {
			fmt.Println(err)
		}
	}()

	numOfJobs := 50000

	for j := 0; j < numOfJobs; j++ {
		randDuration := time.Duration(j) * time.Nanosecond
		//push job
		retries := 1
		go workers.PushJob(uint(j), uint8(retries), func() error {
			//example of heavy task
			time.Sleep(randDuration)

			if j > int(numOfJobs*75/100) {
				return errors.New("job should retry")
			}
			return nil
		})
	}

	// create term so the app didn't exit
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case <-term:
		log.Println("terminate app")
	}

}
