package aws

import "github.com/vinneyto/splatra/api/internal/core"

type Module struct {
	AuthProvider      core.AuthProvider
	JobRepository     core.JobRepository
	ResultURLResolver core.ResultURLResolver
}

func NewModule(_ Config) (*Module, error) {
	return &Module{
		AuthProvider:      &authProviderStub{},
		JobRepository:     &jobRepositoryStub{},
		ResultURLResolver: &resultURLResolverStub{},
	}, nil
}
