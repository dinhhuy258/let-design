package usecase

import (
	"context"
	"job-scheduler-service/internal/entity"
)

type fairSchedulerUsecase struct{}

func NewFairSchedulerUsecase() SchedulerUsecase {
	return &fairSchedulerUsecase{}
}

func (*fairSchedulerUsecase) ScheduleJobs(ctx context.Context, jobs []entity.Job) error {
	return nil
}
