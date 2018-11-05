TKP WORKERS
======

TKP Workers is golang worker thread pool that enforces constant processing rate to do particular job.

Jobs are buffered to channel and spawned workers will consume buffered jobs concurrently.

Behavior :
If job failed, the job will be retried by the workers.

## How to use

1. Create Workers Pool
```go
  //Spawn 50 workers, buffer size = 2056 job channels
  numOfWorkers := uint(50)

  workers := NewWorkers(numOfWorkers)
  workers.Run()
```

2. Spawn Listener
```go
//create a listener to queued jobs
go func() {
      for err := range workers.PollJob() {
        fmt.Println(err)
      }
}()
```

3. Push Job
```go
  numOfJobs := 5000

  for j := 0; j < numOfJobs; j++ {
      randDuration := time.Duration(j) * time.Millisecond
      //push job
      retries := 1
      go workers.PushJob(uint(j), uint8(retries), func() error {
        //example of heavy task
        time.Sleep(randDuration)

        if(j > numOfJobs*0.75){
          return errors.New("should retry")
        }
        return nil
      })
  }

  // create signalTerm so the app didn't exit
  term := make(chan os.Signal, 1)
  signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
  select {
  case <-term:
    log.Println("😥 Signal terminate detected")
  }
```
