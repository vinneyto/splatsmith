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

type ReconstructionJobRepository interface {
	List(ctx context.Context, filter JobListFilter) ([]JobSummary, error)
	GetByID(ctx context.Context, userID, jobID string) (*JobDetails, error)
	FindByIdempotencyKey(ctx context.Context, userID, idempotencyKey string) (*JobDetails, error)
	CreateQueued(ctx context.Context, userID string, req SubmitJobRequest) (*JobDetails, error)
	SetRunning(ctx context.Context, jobID string) error
	SetProgress(ctx context.Context, jobID string, progressPercent int, currentStep string) error
	SetDone(ctx context.Context, jobID string, outputFiles []OutputFileRef) error
	SetFailed(ctx context.Context, jobID, errorMessage string) error
	SetCancelled(ctx context.Context, userID, jobID string) (*JobDetails, error)
	ResetForRetry(ctx context.Context, userID, jobID string) (*JobDetails, error)
}

type ReconstructionSubmissionDispatcher interface {
	Enqueue(ctx context.Context, req ReconstructionSubmissionRequest) error
}

type ReconstructionResultURLResolver interface {
	ResolveResultURL(ctx context.Context, key string, ttl time.Duration) (ResultFileURL, error)
}
