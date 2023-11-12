package entity

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrAuthFailed           = errors.New("authentication failed")
	ErrJobNotFound          = errors.New("job not found")
	ErrJobCannotBeCancelled = errors.New("job cannot be cancelled")
	ErrBadRequest           = errors.New("bad request")
)
