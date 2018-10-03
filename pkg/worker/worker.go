package worker

import (
	"fmt"
	"sync"
)

var (
	once   sync.Once
	wg     sync.WaitGroup
	worker *WorkerImpl
)

type Worker interface {
	Run()
	PushJob(jobID uint, retries uint8, f func() error)
	ConsumeJob(workerID uint, jobs <-chan *Job, errors chan<- error)
	PollJob() chan error
	Close()
}

type WorkerImpl struct {
	wg           *sync.WaitGroup
	maxWorkers   uint
	JobChannel   chan *Job
	ErrorChannel chan error
}

type Job struct {
	ID      uint
	Retries uint8
	F       func() error
}

func NewWorkers(maxWorkers uint, size uint) Worker {
	once.Do(func() {
		worker = &WorkerImpl{
			wg:           &wg,
			maxWorkers:   maxWorkers,
			JobChannel:   make(chan *Job, size),
			ErrorChannel: make(chan error, size),
		}
	})
	return worker
}

func (w WorkerImpl) Run() {
	for i := 0; i < int(w.maxWorkers); i++ {
		go w.ConsumeJob(uint(i), w.JobChannel, w.ErrorChannel)
	}
	w.wg.Wait()
}

func (w WorkerImpl) PushJob(jobID uint, retries uint8, f func() error) {
	select {
	case w.JobChannel <- w.Job(jobID, retries, f):
		w.wg.Add(1)
	default:
		fmt.Printf("dropping job of %v\n", jobID)
	}
}

func (w WorkerImpl) Job(jobID uint, retries uint8, f func() error) *Job {
	return &Job{ID: jobID, Retries: retries, F: func() error {
		if err := f(); err != nil {
			return err
		}
		return nil
	}}
}

func (w WorkerImpl) ConsumeJob(workerID uint, jobs <-chan *Job, errors chan<- error) {
	for job := range jobs {
		// fmt.Println("Exec : JobID", job.ID, "workerID", workerID)
		if err := job.F(); err != nil {
			if job.Retries > 0 {
				// fmt.Println("Retry : JobID", job.ID)
				w.PushJob(job.ID, uint8(job.Retries-1), job.F)
			} else {
				select {
				case errors <- err:
				default:
				}
			}
		}
		w.wg.Done()
	}
}

func (w WorkerImpl) PollJob() chan error {
	return w.ErrorChannel
}

func (w WorkerImpl) Close() {
	close(w.JobChannel)
	close(w.ErrorChannel)
}
