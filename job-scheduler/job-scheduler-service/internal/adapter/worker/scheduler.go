package worker

import (
	"job-scheduler-service/config"
	"job-scheduler-service/pkg/logger"
	"time"
)

type SchedulerWorker interface {
	Start()
	Stop()
}

type schedulerWorker struct {
	Closed chan struct{}
	Ticker *time.Ticker
	// list of shard ids that this worker is responsible for
	ShardIds []uint64
	Logger   *logger.Logger
}

func NewSchedulerWorker(conf *config.Config, logger *logger.Logger) SchedulerWorker {
	return &schedulerWorker{
		Closed:   make(chan struct{}),
		Ticker:   time.NewTicker(time.Second),
		ShardIds: conf.App.ShardIds,
		Logger:   logger,
	}
}

// Start starts the scheduler worker.
func (_self *schedulerWorker) Start() {
	go func() {
		for {
			select {
			case <-_self.Closed:
				return
			case <-_self.Ticker.C:
			}
		}
	}()
}

// Stop stops the scheduler worker.
func (_self *schedulerWorker) Stop() {
	close(_self.Closed)
}
