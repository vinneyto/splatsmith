package core

import (
	"context"
	"time"
)

type AuthRequest struct {
	AuthorizationHeader string
	OIDCIdentityHeader  string
	OIDCDataHeader      string
}

type AuthProvider interface {
	ValidateToken(ctx context.Context, token string) (AuthClaims, error)
}

type AuthRequestAdapter interface {
	Adapt(ctx context.Context, req AuthRequest) (context.Context, AuthRequest)
}

type JobRepository interface {
	List(ctx context.Context, filter JobListFilter) ([]JobSummary, error)
	GetByID(ctx context.Context, userID, jobID string) (*JobDetails, error)
}

type ResultURLResolver interface {
	ResolveResultURL(ctx context.Context, key string, ttl time.Duration) (ResultFileURL, error)
}
