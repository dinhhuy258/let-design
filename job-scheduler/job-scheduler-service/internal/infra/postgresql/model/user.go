package model

import (
	"job-scheduler-service/internal/entity"
)

type User struct {
	BaseModel
	JobWeight float32 `json:"job_weight"`
}

func (User) FromEntity(user entity.User) User {
	return User{
		BaseModel: BaseModel{
			Id: user.Id,
		},
		JobWeight: user.JobWeight,
	}
}

func (_self User) ToEntity() entity.User {
	return entity.User{
		Id:        _self.Id,
		JobWeight: _self.JobWeight,
	}
}

func (User) TableName() string {
	return "users"
}
