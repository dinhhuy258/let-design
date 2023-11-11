package usecase

import (
	"context"
	"job-scheduler-service/internal/entity"
	"job-scheduler-service/pkg/logger"
)

type SchedulerUsecase interface {
	ScheduleJobs(ctx context.Context, jobs []entity.Job) error
}

type baseSchedulerUsecase struct {
	scheduledJobTopic    string
	jobRepository        entity.JobRepository
	messageBusRepository entity.MessageBusRepository
	logger               *logger.Logger
}

func (_self *baseSchedulerUsecase) scheduleJob(ctx context.Context, job entity.Job) error {
	// the following code will cause dual writes problem
	// see: https://developers.redhat.com/articles/2021/07/30/avoiding-dual-writes-event-driven-applications
	// to solve this problem, we can use:
	// - outbox pattern
	// - retry + idempotency
	// for the sake of simplicity, I accept this dual writes problem
	err := _self.messageBusRepository.PublishScheduledEvent(_self.scheduledJobTopic, entity.ScheduledEvent{
		Message: job.Message,
	})
	if err != nil {
		_self.logger.Error("failed to publish scheduled event %v", err)

		return err
	}

	// update job status to running
	job.Status = entity.JobStatusRunning

	err = _self.jobRepository.Update(ctx, job)
	if err != nil {
		// if error occurs, dual writes problem will happen
		// the jobs will be scheduled twice
		// to mitigate this problem, we can use:
		// - outbox pattern
		// - handle idempotency in the executor workers
		_self.logger.Error("failed to update job status to running %v", err)

		return err
	}

	return nil
}
