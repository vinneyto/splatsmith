package core

import "errors"

var (
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInvalidToken   = errors.New("invalid token")
	ErrNotImplemented = errors.New("not implemented")
	ErrJobNotFound    = errors.New("job not found")
)
