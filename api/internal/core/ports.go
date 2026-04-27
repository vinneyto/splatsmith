package core

import (
	"context"
	"io"
)

type AuthProvider interface {
	ValidateToken(ctx context.Context, token string) (AuthClaims, error)
}

type ScanRepository interface {
	Create(ctx context.Context, scan *Scan) error
	GetByID(ctx context.Context, userID, scanID string) (*Scan, error)
	ListByUser(ctx context.Context, userID string, limit, offset int) ([]Scan, error)
	UpdateStatus(
		ctx context.Context,
		scanID string,
		status ScanStatus,
		progressPercent int,
		resultAssetURL *string,
		errorMessage *string,
	) error
}

type ObjectStorage interface {
	SaveInputVideo(ctx context.Context, userID, scanID string, data io.Reader) (string, error)
	OpenResultAsset(ctx context.Context, assetPath string) (io.ReadCloser, error)
}

type PipelineClient interface {
	StartScan(ctx context.Context, scanID, inputVideoPath string) (pipelineJobID string, err error)
}

type Notifier interface {
	NotifyScanCompleted(ctx context.Context, userEmail string, scan Scan) error
}
