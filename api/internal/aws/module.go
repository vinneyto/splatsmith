package aws

import "github.com/vinneyto/splatmaker/api/internal/core"

type Module struct {
	AuthProvider      core.AuthProvider
	LoginProvider     core.LoginProvider
	JobRepository     core.JobRepository
	JobDispatcher     core.JobDispatcher
	ResultURLResolver core.ResultURLResolver
}

func NewModule(_ Config) (*Module, error) {
	return &Module{
		AuthProvider:      &authProviderStub{},
		LoginProvider:     &loginProviderStub{},
		JobRepository:     &jobRepositoryStub{},
		JobDispatcher:     &jobDispatcherStub{},
		ResultURLResolver: &resultURLResolverStub{},
	}, nil
}
