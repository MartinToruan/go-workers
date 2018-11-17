package worker

import (
	"sync"
)

var (
	once   sync.Once
	worker *WorkerImpl
)

type Worker interface {
	Run()
	Coordinator()
	PushJob(jobID uint, retries uint8, f JobFunc)
	ConsumeJob(workerID uint, jobs <-chan *Job, errors chan<- error)
	PollJob() chan error
}

type WorkerImpl struct {
	MaxWorkers   uint
	JobEntry     chan *Job
	JobExec      chan *Job
	ErrorChannel chan error
	Pipeline     []*Job
}

type JobFunc func() error

type Job struct {
	ID         uint
	Retries    uint8
	MaxRetries uint8
	Func       JobFunc
}

func NewWorkers(maxWorkers uint) Worker {
	once.Do(func() {
		worker = &WorkerImpl{
			MaxWorkers:   maxWorkers,
			JobEntry:     make(chan *Job),
			JobExec:      make(chan *Job),
			ErrorChannel: make(chan error),
			Pipeline:     nil,
		}
	})
	return worker
}

func (w WorkerImpl) Run() {
	for i := 0; i < int(w.MaxWorkers); i++ {
		go w.ConsumeJob(uint(i), w.JobExec, w.ErrorChannel)
	}
	go w.Coordinator()
}

func (w WorkerImpl) Coordinator() {
	for {
		select {
		case newJob := <-w.JobEntry:
			var jobFromPipeline *Job
			w.Pipeline = append(w.Pipeline, newJob)
			jobFromPipeline, w.Pipeline = w.Pipeline[0], w.Pipeline[1:]
			go func(job *Job) {
				w.JobExec <- jobFromPipeline
			}(jobFromPipeline)
		}
	}
}

func (w WorkerImpl) PushJob(jobID uint, maxRetries uint8, f JobFunc) {
	go func() {
		w.JobEntry <- w.Job(jobID, maxRetries, maxRetries, f)
	}()
}

func (w WorkerImpl) RetryJob(jobID uint, retries uint8, maxRetries uint8, f JobFunc) {
	go func() {
		w.JobEntry <- w.Job(jobID, retries, maxRetries, f)
	}()
}

func (w WorkerImpl) Job(jobID uint, retries uint8, maxRetries uint8, f JobFunc) *Job {
	return &Job{ID: jobID, Retries: retries, MaxRetries: maxRetries, Func: f}
}

func (w WorkerImpl) ConsumeJob(workerID uint, jobs <-chan *Job, errors chan<- error) {
	for job := range jobs {
		if err := job.Func(); err != nil {
			if job.Retries > 0 {
				go w.RetryJob(job.ID, uint8(job.Retries-1), job.MaxRetries, job.Func)
			} else {
				go func() {
					errors <- err
				}()
			}
		}
	}
}

func (w WorkerImpl) PollJob() chan error {
	return w.ErrorChannel
}
