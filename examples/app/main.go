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
	numOfWorkers := uint(50)
	channelBuffer := uint(2056)

	workers := w.NewWorkers(numOfWorkers, channelBuffer)
	workers.Run()

	//create a listener to queued jobs
	go func() {
		for err := range workers.PollJob() {
			fmt.Println(err)
		}
	}()

	numOfJobs := 5000

	for j := 0; j < numOfJobs; j++ {
		//redeclare to be accessible in the closure
		jobID := j
		retries := 1
		randDuration := time.Duration(jobID) * time.Millisecond

		//job closure
		job := func() error {
			fmt.Println("jobID", jobID, "waits for", randDuration)
			//example of heavy task
			time.Sleep(randDuration)

			if jobID > int(numOfJobs*3/4) {
				return errors.New(fmt.Sprintf("jobID %v more should retry", jobID))
			}
			return nil
		}

		go workers.PushJob(uint(jobID), uint8(retries), job)
	}

	// create term so the app didn't exit
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case <-term:
		log.Println("terminate app")
	}

}
