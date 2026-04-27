package core

import "time"

type ScanStatus string

const (
	ScanStatusInProgress ScanStatus = "in_progress"
	ScanStatusCompleted  ScanStatus = "completed"
	ScanStatusFailed     ScanStatus = "failed"
)

type UserIdentity struct {
	UserID string
	Email  string
}

type AuthClaims struct {
	UserID string
	Email  string
}

type Scan struct {
	ScanID          string
	UserID          string
	Status          ScanStatus
	ProgressPercent int
	InputVideoPath  string
	ResultAssetURL  *string
	PipelineJobID   *string
	ErrorMessage    *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CompletedAt     *time.Time
}
