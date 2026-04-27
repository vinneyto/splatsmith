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

type JobStatus string

const (
	JobStatusNew        JobStatus = "new"
	JobStatusInProgress JobStatus = "in_progress"
	JobStatusDone       JobStatus = "done"
	JobStatusFailed     JobStatus = "failed"
)

type JobListFilter struct {
	UserID string
	Status *JobStatus
	Limit  int
	Offset int
}

type JobSummary struct {
	JobID     string
	UserID    string
	Status    JobStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OutputFileRef struct {
	Key       string
	FileName  string
	SizeBytes *int64
}

type JobDetails struct {
	Summary      JobSummary
	ErrorMessage *string
	OutputFiles  []OutputFileRef
}

type ResultFileURL struct {
	Key       string
	FileName  string
	URL       string
	ExpiresAt time.Time
}
