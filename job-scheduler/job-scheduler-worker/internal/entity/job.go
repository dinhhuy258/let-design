package entity

import (
	"context"
	"time"
)

const (
	JobStatusCreated = "created"
	JobStatusRunning = "running"
)

type Job struct {
	Id        uint64
	UserId    uint64
	Message   string
	Status    string
	ExecuteAt time.Time
	ShardId   uint64
	User      *User
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
