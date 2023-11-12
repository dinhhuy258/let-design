package usecase

import (
	"context"
	"job-service/internal/entity"
	"job-service/pkg/logger"
)

type JobUsecase interface {
	GetJobs(ctx context.Context, userId uint64) ([]entity.Job, error)
	CreateJob(ctx context.Context, job entity.Job) (entity.Job, error)
	CancelJob(ctx context.Context, userId, jobId uint64) error
	GetJob(ctx context.Context, userId, jobId uint64) (entity.Job, error)
}

type jobUsecase struct {
	jobShardUsecase JobShardUsecase
	jobRepository   entity.JobRepository
	logger          *logger.Logger
}

func NewJobUsecase(
	jobShardUsecase JobShardUsecase,
	jobRepository entity.JobRepository,
	logger *logger.Logger,
) JobUsecase {
	return &jobUsecase{
		jobShardUsecase: jobShardUsecase,
		jobRepository:   jobRepository,
		logger:          logger,
	}
}

func (_self *jobUsecase) CreateJob(ctx context.Context, job entity.Job) (entity.Job, error) {
	_self.logger.Info("creating job %v", job)

	job.Status = entity.JobStatusCreated
	job.ShardId = _self.jobShardUsecase.GetSharedId()

	job, err := _self.jobRepository.Create(ctx, job)
	if err != nil {
		_self.logger.Error("failed to create job %v", err)

		return entity.Job{}, err
	}

	return job, nil
}

func (_self *jobUsecase) CancelJob(ctx context.Context, userId, jobId uint64) error {
	_self.logger.Info("updating job id %d", jobId)

	job, err := _self.jobRepository.FindById(ctx, jobId)
	if err != nil {
		_self.logger.Error("failed to find job %v", err)

		return err
	}

	if job == nil {
		return entity.ErrJobNotFound
	}

	if job.UserId != userId {
		return entity.ErrJobNotFound
	}

	// FIXME: This line of code is not work perfectly in concurrent environment
	// as I just only focus on the scheduler part, I will leave it like this for now
	if job.Status != entity.JobStatusCreated {
		// Only job with status created can be cancelled
		return entity.ErrJobCannotBeCancelled
	}

	job.Status = entity.JobStatusCancelled

	err = _self.jobRepository.Update(ctx, *job)
	if err != nil {
		_self.logger.Error("failed to update job id %d. Error %v", jobId, err)
	}

	return err
}

func (_self *jobUsecase) GetJobs(ctx context.Context, userId uint64) ([]entity.Job, error) {
	jobs, err := _self.jobRepository.FindMultiByUserId(ctx, userId)
	if err != nil {
		_self.logger.Error("failed to get jobs %v", err)

		return nil, err
	}

	return jobs, nil
}

func (_self *jobUsecase) GetJob(ctx context.Context, userId uint64, jobId uint64) (entity.Job, error) {
	job, err := _self.jobRepository.FindById(ctx, jobId)
	if err != nil {
		_self.logger.Error("failed to get job %v", err)

		return entity.Job{}, err
	}

	if job == nil {
		return entity.Job{}, entity.ErrJobNotFound
	}

	if job.UserId != userId {
		return entity.Job{}, entity.ErrJobNotFound
	}

	return *job, nil
}
