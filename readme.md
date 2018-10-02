TKP WORKERS
======

TKP Workers is golang worker thread pool that enforces constant processing rate to do particular job.

Jobs are buffered to channel and spawned workers will consume buffered jobs concurrently.

Behavior :
1. If the channel buffer is full, the job will be disbanded.
2. If job failed, the job will be retried by the workers.

See examples to use this with **TDK - REST/ GRPC**

## How to use

1. Create Workers Pool
```go
  //Spawn 50 workers, buffer size = 2056 job channels
  numOfWorkers := 50
  channelBuffer := 2056

  workers := NewWorkers(numOfWorkers, channelBuffer)
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
        return nil
      })
  }
  //forever loop
  for { }
```
