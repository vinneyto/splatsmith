package core

import "errors"

var (
	ErrUnauthorized           = errors.New("unauthorized")
	ErrInvalidToken           = errors.New("invalid token")
	ErrNotImplemented         = errors.New("not implemented")
	ErrJobNotFound            = errors.New("job not found")
	ErrInvalidArgument        = errors.New("invalid argument")
	ErrConflict               = errors.New("conflict")
	ErrIdempotencyKeyRequired = errors.New("idempotency key is required")
	ErrJobNotCancelable       = errors.New("job is not cancelable")
	ErrJobNotRetryable        = errors.New("job is not retryable")
)
