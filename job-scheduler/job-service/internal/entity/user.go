package entity

import "context"

type User struct {
	Id       uint64 `json:"id"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserRepository interface {
	Create(ctx context.Context, user User) (User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
}
