package core

import "errors"

var (
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInvalidToken      = errors.New("invalid token")
	ErrNotImplemented    = errors.New("not implemented")
	ErrScanNotFound      = errors.New("scan not found")
	ErrInvalidScanStatus = errors.New("invalid scan status")
)
