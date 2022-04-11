package runner

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

var workersNumber int64 = 3
var drainTimeout = time.Second * 5

type ExecutionWorkerPool struct {
	runnerService RunnerService
}

func NewExecutionWorkerPool(runnerService RunnerService) *ExecutionWorkerPool {
	return &ExecutionWorkerPool{
		runnerService: runnerService,
	}
}

// Run runs a pool of workers to process the execution requests
func (e *ExecutionWorkerPool) Run(ctx context.Context) {
	log.Infof("Starting execution pool. Workers limit: %d", workersNumber)
	sem := semaphore.NewWeighted(workersNumber)
	channel := e.runnerService.GetChannel()

	for {
		select {
		case execution := <-channel:
			if err := sem.Acquire(ctx, 1); err != nil {
				log.Debugf("Discarding execution: %d, shutting down already.", execution.ExecutionID)
				break
			}

			go func() {
				defer sem.Release(1)
				e.runnerService.Execute(execution)
			}()
		case <-ctx.Done():
			log.Infof("Projectors worker pool is shutting down... Waiting for active workers to drain.")

			ctx, cancel := context.WithTimeout(context.Background(), drainTimeout)
			defer cancel()

			if err := sem.Acquire(ctx, workersNumber); err != nil {
				log.Warnf("Timed out while draining workers: %v", err)
			}

			return
		}
	}
}
