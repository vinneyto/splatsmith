package aws

import (
	"context"
	"fmt"
	"io"

	"github.com/vinneyto/ariadne/api/internal/core"
)

type authProviderStub struct{}

type scanRepositoryStub struct{}

type objectStorageStub struct{}

type pipelineClientStub struct{}

type notifierStub struct{}

func (s *authProviderStub) ValidateToken(context.Context, string) (core.AuthClaims, error) {
	return core.AuthClaims{}, fmt.Errorf("aws auth provider: %w", core.ErrNotImplemented)
}

func (s *scanRepositoryStub) Create(context.Context, *core.Scan) error {
	return fmt.Errorf("aws scan repository: %w", core.ErrNotImplemented)
}

func (s *scanRepositoryStub) GetByID(context.Context, string, string) (*core.Scan, error) {
	return nil, fmt.Errorf("aws scan repository: %w", core.ErrNotImplemented)
}

func (s *scanRepositoryStub) ListByUser(context.Context, string, int, int) ([]core.Scan, error) {
	return nil, fmt.Errorf("aws scan repository: %w", core.ErrNotImplemented)
}

func (s *scanRepositoryStub) UpdateStatus(context.Context, string, core.ScanStatus, int, *string, *string) error {
	return fmt.Errorf("aws scan repository: %w", core.ErrNotImplemented)
}

func (s *objectStorageStub) SaveInputVideo(context.Context, string, string, io.Reader) (string, error) {
	return "", fmt.Errorf("aws object storage: %w", core.ErrNotImplemented)
}

func (s *objectStorageStub) OpenResultAsset(context.Context, string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("aws object storage: %w", core.ErrNotImplemented)
}

func (s *pipelineClientStub) StartScan(context.Context, string, string) (string, error) {
	return "", fmt.Errorf("aws pipeline client: %w", core.ErrNotImplemented)
}

func (s *notifierStub) NotifyScanCompleted(context.Context, string, core.Scan) error {
	return fmt.Errorf("aws notifier: %w", core.ErrNotImplemented)
}
