package model

import (
	"job-service/internal/entity"
	"time"
)

type Job struct {
	BaseModel
	UserId       uint64    `gorm:"user_id"`
	ShardId      uint64    `gorm:"shard_id"`
	Message      string    `gorm:"message"`
	Status       string    `gorm:"status"`
	WeightFactor float32   `gorm:"weight_factor"`
	ExecuteAt    time.Time `gorm:"execute_at"`
}

func (Job) FromEntity(job entity.Job) Job {
	return Job{
		BaseModel: BaseModel{
			Id: job.Id,
		},
		UserId:       job.UserId,
		ShardId:      job.ShardId,
		Message:      job.Message,
		Status:       job.Status,
		WeightFactor: job.WeightFactor,
		ExecuteAt:    job.ExecuteAt,
	}
}

func (_self Job) ToEntity() entity.Job {
	return entity.Job{
		Id:           _self.Id,
		UserId:       _self.UserId,
		ShardId:      _self.ShardId,
		Message:      _self.Message,
		Status:       _self.Status,
		WeightFactor: _self.WeightFactor,
		ExecuteAt:    _self.ExecuteAt,
	}
}

func (Job) TableName() string {
	return "jobs"
}
