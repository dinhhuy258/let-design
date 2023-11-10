package worker

import (
	"job-scheduler-service/config"
	"job-scheduler-service/internal/usecase"
	"job-scheduler-service/pkg/logger"
	"time"
)

type SchedulerWorker interface {
	Start()
	Stop()
}

type schedulerWorker struct {
	closed chan struct{}
	ticker *time.Ticker
	// list of shard ids that this worker is responsible for
	shardIds []uint64

	jobUsecase       usecase.JobUsecase
	schedulerUsecase usecase.SchedulerUsecase

	logger *logger.Logger
}

func NewSchedulerWorker(
	jobUsecase usecase.JobUsecase,
	schedulerUsecase usecase.SchedulerUsecase,
	conf *config.Config,
	logger *logger.Logger,
) SchedulerWorker {
	return &schedulerWorker{
		closed:           make(chan struct{}),
		ticker:           time.NewTicker(time.Second),
		shardIds:         conf.App.ShardIds,
		jobUsecase:       jobUsecase,
		schedulerUsecase: schedulerUsecase,
		logger:           logger,
	}
}

// Start starts the scheduler worker.
func (_self *schedulerWorker) Start() {
	go func() {
		for {
			select {
			case <-_self.closed:
				return
			case <-_self.ticker.C:
			}
		}
	}()
}

// Stop stops the scheduler worker.
func (_self *schedulerWorker) Stop() {
	close(_self.closed)
}
