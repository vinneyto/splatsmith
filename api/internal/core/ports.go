package core

import (
	"context"
	"time"
)

type AuthProvider interface {
	ValidateToken(ctx context.Context, token string) (AuthClaims, error)
}

type LoginProvider interface {
	LoginWithPassword(ctx context.Context, username, password string) (LoginResult, error)
}

type JobRepository interface {
	List(ctx context.Context, filter JobListFilter) ([]JobSummary, error)
	GetByID(ctx context.Context, userID, jobID string) (*JobDetails, error)
}

type ResultURLResolver interface {
	ResolveResultURL(ctx context.Context, key string, ttl time.Duration) (ResultFileURL, error)
}
