# ninhydrin-go

Ninhydrin API Go Client

```go
package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/lostinsoba/ninhydrin-go"
)

func main() {

	api := ninhydrin.New()
	ctx := context.Background()

	tasks, err := api.Task.Capture(ctx, "Infrastructure", 1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, task := range tasks {
		switch task.ID {
		case "LongRunningTask":
			var taskStatus ninhydrin.TaskStatus
			taskErr := performTask(longRunningTask, task.Timeout)
			if taskErr != nil {
				taskStatus = ninhydrin.TaskStatusFailed
			} else {
				taskStatus = ninhydrin.TaskStatusDone
			}
			releaseErr := api.Task.Release(ctx, taskStatus, []string{taskID})
			if releaseErr != nil {
				fmt.Println(releaseErr)
			}
		}
	}
}

func performTask(task func() error, timeoutSeconds int64) error {
	var (
		timeoutDuration = time.Duration(timeoutSeconds) * time.Second
		timeoutTicker   = time.NewTicker(timeoutDuration)
		errChan         = make(chan error)
	)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		errChan <- task()
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var err error
	select {
	case err = <-errChan:
		break
	case <-timeoutTicker.C:
		fmt.Println("timeout")
		os.Exit(1)
	}

	return err
}

func longRunningTask() error {
	time.Sleep(5 * time.Second)
	return nil
}
```