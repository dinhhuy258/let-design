package usecase

import (
	"context"
	"job-scheduler-service/internal/entity"
)

type SchedulerUsecase interface {
	ScheduleJobs(ctx context.Context, jobs []entity.Job) error
}

type baseSchedulerUsecase struct{}

func (_self *baseSchedulerUsecase) scheduleJob(ctx context.Context, job entity.Job) error {
	// push job to kafka
	print(job.Message)

	return nil
}
