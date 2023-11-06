package entity

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrAuthFailed   = errors.New("authentication failed")
)
