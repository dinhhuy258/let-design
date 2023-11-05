package entity

import "context"

type User struct {
	Id       uint64
	Username string
	Password string
}

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
}
