package core

import (
	"context"
	"time"
)

type JobRepository interface {
	List(ctx context.Context, filter JobListFilter) ([]JobSummary, error)
	GetByID(ctx context.Context, jobID string) (*JobDetails, error)
}

type ResultURLResolver interface {
	ResolveResultURL(ctx context.Context, key string, ttl time.Duration) (ResultFileURL, error)
}
