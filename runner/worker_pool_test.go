package runner

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WorkerPoolTestCase struct {
	suite.Suite
}

func TestWorkerPoolTestCase(t *testing.T) {
	suite.Run(t, new(WorkerPoolTestCase))
}

func (suite *WorkerPoolTestCase) Test_Run() {
	channel := make(chan *ExecutionEvent)
	execution := &ExecutionEvent{ExecutionID: uuid.New()}

	var wg sync.WaitGroup
	wg.Add(2)

	mockRunnerService := new(MockRunnerService)
	mockRunnerService.On("Execute", mock.Anything).Run(func(args mock.Arguments) {
		wg.Done()
	}).Return(nil)
	mockRunnerService.On("GetChannel").Return(channel)

	workerPool := NewExecutionWorkerPool(mockRunnerService)

	ctx, cancel := context.WithCancel(context.Background())

	go workerPool.Run(ctx)
	channel <- execution
	channel <- execution

	wg.Wait()

	mockRunnerService.AssertNumberOfCalls(suite.T(), "Execute", 2)
	cancel()
}
