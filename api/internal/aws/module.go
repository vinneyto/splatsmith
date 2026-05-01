package aws

import "github.com/vinneyto/splatmaker/api/internal/core"

type Module struct {
	AuthProvider                       core.AuthProvider
	LoginProvider                      core.LoginProvider
	ReconstructionJobRepository        core.ReconstructionJobRepository
	ReconstructionSubmissionDispatcher core.ReconstructionSubmissionDispatcher
	ResultURLResolver                  core.ResultURLResolver
}

func NewModule(_ Config) (*Module, error) {
	return &Module{
		AuthProvider:                       &authProviderStub{},
		LoginProvider:                      &loginProviderStub{},
		ReconstructionJobRepository:        &jobRepositoryStub{},
		ReconstructionSubmissionDispatcher: &jobDispatcherStub{},
		ResultURLResolver:                  &resultURLResolverStub{},
	}, nil
}
