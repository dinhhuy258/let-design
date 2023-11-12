package usecase

import (
	"context"
	"job-scheduler-worker/internal/entity"
	"time"
)

type JobUsecase interface {
	GetAvailableJobs(ctx context.Context, shardIds []uint64) ([]entity.Job, error)
}

type jobUsecase struct {
	jobRepository entity.JobRepository
}

func NewJobUsecase(jobRepository entity.JobRepository) JobUsecase {
	return &jobUsecase{
		jobRepository: jobRepository,
	}
}

func (_self *jobUsecase) GetAvailableJobs(ctx context.Context, shardIds []uint64) ([]entity.Job, error) {
	now := time.Now()

	return _self.jobRepository.FindMultiAvailableJobs(ctx, entity.JobStatusCreated, now, shardIds)
}
