package entity

import "context"

type User struct {
	Id       uint64
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRepository interface {
	Create(ctx context.Context, user User) (User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
}
