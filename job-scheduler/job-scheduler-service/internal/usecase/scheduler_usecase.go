package usecase

import (
	"context"
	"job-scheduler-service/internal/entity"
)

type SchedulerUsecase interface {
	ScheduleJobs(ctx context.Context, jobs []entity.Job) error
}
