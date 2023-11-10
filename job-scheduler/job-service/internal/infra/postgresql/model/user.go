package model

import (
	"job-service/internal/entity"
)

type User struct {
	BaseModel
	Username  string  `gorm:"username"`
	Password  string  `gorm:"password"`
	JobWeight float32 `json:"job_weight"`
}

func (User) FromEntity(user entity.User) User {
	return User{
		BaseModel: BaseModel{
			Id: user.Id,
		},
		Username:  user.Username,
		Password:  user.Password,
		JobWeight: user.JobWeight,
	}
}

func (_self User) ToEntity() entity.User {
	return entity.User{
		Id:        _self.Id,
		Username:  _self.Username,
		Password:  _self.Password,
		JobWeight: _self.JobWeight,
	}
}

func (User) TableName() string {
	return "users"
}
