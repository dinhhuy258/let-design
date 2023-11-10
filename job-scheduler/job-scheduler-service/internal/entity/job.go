package entity

import (
	"context"
	"time"
)

const (
	JobStatusCreated = "created"
)

type Job struct {
	Id           uint64
	UserId       uint64
	Message      string
	Status       string
	WeightFactor float32
	ExecuteAt    time.Time
	ShardId      uint64
}

type JobRepository interface {
	FindMultiAvailableJobs(
		ctx context.Context,
		status string,
		executeAt time.Time,
		shardIds []uint64,
	) ([]Job, error)
	Update(ctx context.Context, job Job) error
}
