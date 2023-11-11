package usecase

import (
	"context"
	"job-scheduler-service/internal/entity"
)

type SchedulerUsecase interface {
	ScheduleJobs(ctx context.Context, jobs []entity.Job) error
}

type baseSchedulerUsecase struct {
	scheduledJobTopic    string
	messageBusRepository entity.MessageBusRepository
}

func (_self *baseSchedulerUsecase) scheduleJob(job entity.Job) error {
	_self.messageBusRepository.PublishScheduledEvent(_self.scheduledJobTopic, entity.ScheduledEvent{
		Message: job.Message,
	})

	return nil
}
