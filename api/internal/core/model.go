package core

import "time"

type UserIdentity struct {
	UserID string
	Email  string
}

type AuthClaims struct {
	UserID string
	Email  string
}

type LoginResult struct {
	Token string
	User  UserIdentity
}

type JobStatus string

const (
	JobStatusNew        JobStatus = "new"
	JobStatusQueued     JobStatus = "queued"
	JobStatusInProgress JobStatus = "in_progress"
	JobStatusDone       JobStatus = "done"
	JobStatusFailed     JobStatus = "failed"
	JobStatusCancelled  JobStatus = "cancelled"
)

type JobListFilter struct {
	UserID string
	Status *JobStatus
	Limit  int
	Offset int
}

type JobSummary struct {
	JobID           string
	UserID          string
	Status          JobStatus
	ProgressPercent int
	CurrentStep     *string
	IdempotencyKey  *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OutputFileRef struct {
	Key       string
	FileName  string
	SizeBytes *int64
}

type JobDetails struct {
	Summary         JobSummary
	ErrorMessage    *string
	OutputFiles     []OutputFileRef
	Attempt         int
	SourceRef       *string
	SimulateFailure bool
	StartedAt       *time.Time
	FinishedAt      *time.Time
	LastHeartbeatAt *time.Time
}

type ResultFileURL struct {
	Key       string
	FileName  string
	URL       string
	ExpiresAt time.Time
}

type SubmitJobRequest struct {
	IdempotencyKey  string
	Name            *string
	SourceRef       *string
	SimulateFailure bool
}

type SubmitJobResult struct {
	Job     JobDetails
	Created bool
}

type JobDispatchRequest struct {
	JobID           string
	UserID          string
	SimulateFailure bool
	IdempotencyKey  string
	CurrentAttempt  int
}
