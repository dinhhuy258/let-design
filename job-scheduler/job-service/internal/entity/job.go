package entity

import (
	"context"
	"time"
)

const (
	JobStatusCreated   = "created"
	JobStatusRunning   = "running"
	JobStatusFailed    = "failed"
	JobStatusCancelled = "cancelled"
	JobStatusCompleted = "completed"
)

type Job struct {
	Id           uint64    `json:"id"`
	UserId       uint64    `json:"user_id"`
	Message      string    `json:"message"`
	Status       string    `json:"status"`
	ExecuteAt    time.Time `json:"execute_at"`
	ShardId      uint64
}

type JobRepository interface {
	Create(ctx context.Context, job Job) (Job, error)
	FindById(ctx context.Context, id uint64) (*Job, error)
	FindMultiByUserId(ctx context.Context, userId uint64) ([]Job, error)
	Update(ctx context.Context, job Job) error
}
