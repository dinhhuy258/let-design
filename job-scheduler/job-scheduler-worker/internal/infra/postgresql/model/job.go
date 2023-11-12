package model

import (
	"job-scheduler-worker/internal/entity"
	"time"
)

type Job struct {
	BaseModel
	UserId    uint64    `gorm:"user_id"`
	ShardId   uint64    `gorm:"shard_id"`
	Message   string    `gorm:"message"`
	Status    string    `gorm:"status"`
	ExecuteAt time.Time `gorm:"execute_at"`
	User      *User     `gorm:"foreignKey:UserId;references:Id"`
}

func (Job) FromEntity(job entity.Job) Job {
	return Job{
		BaseModel: BaseModel{
			Id: job.Id,
		},
		UserId:    job.UserId,
		ShardId:   job.ShardId,
		Message:   job.Message,
		Status:    job.Status,
		ExecuteAt: job.ExecuteAt,
	}
}

func (_self Job) ToEntity() entity.Job {
	job := entity.Job{
		Id:        _self.Id,
		UserId:    _self.UserId,
		ShardId:   _self.ShardId,
		Message:   _self.Message,
		Status:    _self.Status,
		ExecuteAt: _self.ExecuteAt,
	}

	if _self.User != nil {
		user := _self.User.ToEntity()
		job.User = &user
	}

	return job
}

func (Job) TableName() string {
	return "jobs"
}
