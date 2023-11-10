package entity

import "context"

type User struct {
	Id        uint64  `json:"id"`
	Username  string  `json:"username"`
	Password  string  `json:"password,omitempty"`
	JobWeight float32 `json:"job_weight"`
}

type UserRepository interface {
	Create(ctx context.Context, user User) (User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
}
