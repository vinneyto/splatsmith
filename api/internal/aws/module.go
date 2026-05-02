package aws

import "github.com/vinneyto/splatmaker/api/internal/core"

type Module struct {
	AuthProvider               core.AuthProvider
	LoginProvider              core.LoginProvider
	JobRepository              core.JobRepository
	PipelineSettingsRepository core.PipelineSettingsRepository
	JobDispatcher              core.JobDispatcher
	ResultURLResolver          core.ResultURLResolver
}

func NewModule(_ Config) (*Module, error) {
	return &Module{
		AuthProvider:               &authProviderStub{},
		LoginProvider:              &loginProviderStub{},
		JobRepository:              &jobRepositoryStub{},
		PipelineSettingsRepository: &pipelineSettingsRepositoryStub{},
		JobDispatcher:              &jobDispatcherStub{},
		ResultURLResolver:          &resultURLResolverStub{},
	}, nil
}
