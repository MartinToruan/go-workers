package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	w "workers/pkg/worker"
)

func main() {
	numOfWorkers := uint(10)

	workers := w.NewWorkers(numOfWorkers)
	workers.Run()

	//create a listener to queued jobs
	go func() {
		for err := range workers.PollJob() {
			fmt.Println(err)
		}
	}()

	numOfJobs := 1000000

	for j := 0; j < numOfJobs; j++ {
		//redeclare to be accessible in the closure
		jobID := j
		retries := 2
		// randDuration := time.Duration(jobID) * time.Millisecond

		//job closure
		job := func() error {
			//example of heavy task
			// time.Sleep(randDuration)
			fmt.Printf("job : %v execute \n", jobID)
			//example of job retries
			if jobID > int(numOfJobs*3/4) {
				return errors.New(fmt.Sprintf("job : %v failed to execute", jobID))
			}
			return nil
		}

		workers.PushJob(uint(jobID), uint8(retries), job)
	}
	// create term so the app didn't exit
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case <-term:
		log.Println("terminate app")
	}

}
