package worker

import (
	"context"
	"job-scheduler-worker/config"
	"job-scheduler-worker/internal/usecase"
	"job-scheduler-worker/pkg/logger"
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
		ctx := context.Background()

		for {
			select {
			case <-_self.closed:
				return
			case <-_self.ticker.C:
				jobs, err := _self.jobUsecase.GetAvailableJobs(ctx, _self.shardIds)
				if err != nil {
					_self.logger.Error("error while getting available jobs", err)
				} else {
					err := _self.schedulerUsecase.ScheduleJobs(ctx, jobs)
					if err != nil {
						_self.logger.Error("error while scheduling jobs %v", err)
					}
				}
			}
		}
	}()
}

// Stop stops the scheduler worker.
func (_self *schedulerWorker) Stop() {
	close(_self.closed)
}
