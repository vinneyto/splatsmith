package core

import (
	"context"
	"time"
)

type AuthProvider interface {
	ValidateHeaders(ctx context.Context, headers map[string]string) (AuthClaims, error)
}

type JobRepository interface {
	List(ctx context.Context, filter JobListFilter) ([]JobSummary, error)
	GetByID(ctx context.Context, userID, jobID string) (*JobDetails, error)
}

type ResultURLResolver interface {
	ResolveResultURL(ctx context.Context, key string, ttl time.Duration) (ResultFileURL, error)
}
