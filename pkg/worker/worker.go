package worker

import (
	"sync"
)

var (
	once   sync.Once
	wg     sync.WaitGroup
	worker *WorkerImpl
)

type Worker interface {
	Run()
	Coordinator()
	PushJob(jobID uint, retries uint8, f JobFunc)
	ConsumeJob(workerID uint, jobs <-chan *Job, errors chan<- error)
	PollJob() chan error
	Close()
}

type WorkerImpl struct {
	WG           *sync.WaitGroup
	MaxWorkers   uint
	JobEntry     chan *Job
	JobExec      chan *Job
	ErrorChannel chan error
	Pipeline     []*Job
}

type JobFunc func() error

type Job struct {
	ID      uint
	Retries uint8
	Func    JobFunc
}

func NewWorkers(maxWorkers uint) Worker {
	once.Do(func() {
		worker = &WorkerImpl{
			WG:           &wg,
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
	w.WG.Wait()
}

func (w WorkerImpl) Coordinator() {
	for {
		select {
		case newJob := <-w.JobEntry:
			var jobFromPipeline *Job
			w.Pipeline = append(w.Pipeline, newJob)
			jobFromPipeline, w.Pipeline = w.Pipeline[0], w.Pipeline[1:]
			w.JobExec <- jobFromPipeline
		}
	}
}

func (w WorkerImpl) PushJob(jobID uint, retries uint8, f JobFunc) {
	w.WG.Add(1)
	w.JobEntry <- w.Job(jobID, retries, f)
}

func (w WorkerImpl) Job(jobID uint, retries uint8, f JobFunc) *Job {
	return &Job{ID: jobID, Retries: retries, Func: func() error {
		if err := f(); err != nil {
			return err
		}
		return nil
	}}
}

func (w WorkerImpl) ConsumeJob(workerID uint, jobs <-chan *Job, errors chan<- error) {
	for job := range jobs {
		if err := job.Func(); err != nil {
			if job.Retries > 0 {
				go w.PushJob(job.ID, uint8(job.Retries-1), job.Func)
			} else {
				errors <- err
			}
		}
		w.WG.Done()
	}
}

func (w WorkerImpl) PollJob() chan error {
	return w.ErrorChannel
}

func (w WorkerImpl) Close() {
	close(w.JobEntry)
	close(w.JobExec)
	close(w.ErrorChannel)
}
