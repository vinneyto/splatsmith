package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/vinneyto/splatmaker/api/internal/core"
)

type authProviderStub struct{}
type loginProviderStub struct{}
type jobRepositoryStub struct{}
type pipelineSettingsRepositoryStub struct{}
type jobDispatcherStub struct{}
type resultURLResolverStub struct{}

func (s *authProviderStub) ValidateToken(context.Context, string) (core.AuthClaims, error) {
	return core.AuthClaims{}, fmt.Errorf("aws auth provider: %w", core.ErrNotImplemented)
}

func (s *loginProviderStub) LoginWithPassword(context.Context, string, string) (core.LoginResult, error) {
	return core.LoginResult{}, fmt.Errorf("aws login provider: %w", core.ErrNotImplemented)
}

func (s *jobRepositoryStub) List(context.Context, core.JobListFilter) ([]core.JobSummary, error) {
	return nil, fmt.Errorf("aws job repository: %w", core.ErrNotImplemented)
}

func (s *jobRepositoryStub) GetByID(context.Context, string, string) (*core.JobDetails, error) {
	return nil, fmt.Errorf("aws job repository: %w", core.ErrNotImplemented)
}

func (s *jobRepositoryStub) FindByIdempotencyKey(context.Context, string, string) (*core.JobDetails, error) {
	return nil, fmt.Errorf("aws job repository: %w", core.ErrNotImplemented)
}

func (s *jobRepositoryStub) CreateQueued(context.Context, string, core.SubmitJobRequest) (*core.JobDetails, error) {
	return nil, fmt.Errorf("aws job repository: %w", core.ErrNotImplemented)
}

func (s *jobRepositoryStub) SetRunning(context.Context, string) error {
	return fmt.Errorf("aws job repository: %w", core.ErrNotImplemented)
}

func (s *jobRepositoryStub) SetProgress(context.Context, string, int, string) error {
	return fmt.Errorf("aws job repository: %w", core.ErrNotImplemented)
}

func (s *jobRepositoryStub) SetDone(context.Context, string, []core.OutputFileRef) error {
	return fmt.Errorf("aws job repository: %w", core.ErrNotImplemented)
}

func (s *jobRepositoryStub) SetFailed(context.Context, string, string) error {
	return fmt.Errorf("aws job repository: %w", core.ErrNotImplemented)
}

func (s *jobRepositoryStub) SetCancelled(context.Context, string, string) (*core.JobDetails, error) {
	return nil, fmt.Errorf("aws job repository: %w", core.ErrNotImplemented)
}

func (s *jobRepositoryStub) ResetForRetry(context.Context, string, string) (*core.JobDetails, error) {
	return nil, fmt.Errorf("aws job repository: %w", core.ErrNotImplemented)
}

func (s *jobDispatcherStub) Enqueue(context.Context, core.JobDispatchRequest) error {
	return fmt.Errorf("aws dispatcher: %w", core.ErrNotImplemented)
}

func (s *pipelineSettingsRepositoryStub) List(context.Context, core.PipelineSettingsListFilter) ([]core.PipelineSettingsRecord, error) {
	return nil, fmt.Errorf("aws pipeline settings repository: %w", core.ErrNotImplemented)
}

func (s *pipelineSettingsRepositoryStub) GetByID(context.Context, string, string) (*core.PipelineSettingsRecord, error) {
	return nil, fmt.Errorf("aws pipeline settings repository: %w", core.ErrNotImplemented)
}

func (s *pipelineSettingsRepositoryStub) Create(context.Context, core.CreatePipelineSettingsInput) (*core.PipelineSettingsRecord, error) {
	return nil, fmt.Errorf("aws pipeline settings repository: %w", core.ErrNotImplemented)
}

func (s *pipelineSettingsRepositoryStub) Update(context.Context, core.UpdatePipelineSettingsInput) (*core.PipelineSettingsRecord, error) {
	return nil, fmt.Errorf("aws pipeline settings repository: %w", core.ErrNotImplemented)
}

func (s *pipelineSettingsRepositoryStub) Delete(context.Context, string, string) error {
	return fmt.Errorf("aws pipeline settings repository: %w", core.ErrNotImplemented)
}

func (s *resultURLResolverStub) ResolveResultURL(context.Context, string, time.Duration) (core.ResultFileURL, error) {
	return core.ResultFileURL{}, fmt.Errorf("aws result url resolver: %w", core.ErrNotImplemented)
}
