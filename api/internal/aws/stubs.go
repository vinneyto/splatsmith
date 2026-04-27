package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/vinneyto/splatra/api/internal/core"
)

type authProviderStub struct{}
type jobRepositoryStub struct{}
type resultURLResolverStub struct{}

func (s *authProviderStub) ValidateToken(context.Context, string) (core.AuthClaims, error) {
	return core.AuthClaims{}, fmt.Errorf("aws auth provider: %w", core.ErrNotImplemented)
}

func (s *jobRepositoryStub) List(context.Context, core.JobListFilter) ([]core.JobSummary, error) {
	return nil, fmt.Errorf("aws job repository: %w", core.ErrNotImplemented)
}

func (s *jobRepositoryStub) GetByID(context.Context, string, string) (*core.JobDetails, error) {
	return nil, fmt.Errorf("aws job repository: %w", core.ErrNotImplemented)
}

func (s *resultURLResolverStub) ResolveResultURL(context.Context, string, time.Duration) (core.ResultFileURL, error) {
	return core.ResultFileURL{}, fmt.Errorf("aws result url resolver: %w", core.ErrNotImplemented)
}
