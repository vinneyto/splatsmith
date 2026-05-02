package aws

import "github.com/vinneyto/splatmaker/api/internal/core"

type Module struct {
	AuthProvider      core.AuthProvider
	JobRepository     core.JobRepository
	ResultURLResolver core.ResultURLResolver
}

func NewModule(cfg Config) (*Module, error) {
	repo, err := NewDynamoJobRepository(cfg)
	if err != nil {
		return nil, err
	}
	resolver, err := NewS3ResultURLResolver(cfg)
	if err != nil {
		return nil, err
	}
	return &Module{
		AuthProvider:      NewALBAuthProvider(),
		JobRepository:     repo,
		ResultURLResolver: resolver,
	}, nil
}
