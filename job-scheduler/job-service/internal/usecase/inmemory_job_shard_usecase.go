package usecase

import (
	"job-service/config"
	"sync/atomic"
)

// inmemoryJobShardUsecase is an in-memory implementation of JobShardUsecase
// this implementation is just for demo purpose
type inmemoryJobShardUsecase struct {
	shardSize uint64
	jobCount  atomic.Uint64
}

// NewInmemoryJobShardUsecase creates a new instance of inmemoryJobShardUsecase
func NewInmemoryJobShardUsecase(conf *config.Config) JobShardUsecase {
	return &inmemoryJobShardUsecase{
		shardSize: conf.App.JobShardSize,
		jobCount:  atomic.Uint64{},
	}
}

func (_self *inmemoryJobShardUsecase) GetSharedId() uint64 {
	_self.jobCount.Add(1)
	shardId := _self.jobCount.Load() % _self.shardSize

	return shardId
}
