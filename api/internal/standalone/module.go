package standalone

import "github.com/vinneyto/splatmaker/api/internal/core"

type Module struct {
	AuthProvider      core.AuthProvider
	JobRepository     core.JobRepository
	ResultURLResolver core.ResultURLResolver

	closers []func() error
}

func NewModule(cfg Config) (*Module, error) {
	repo, err := NewSQLiteJobRepository(cfg.SQLitePath)
	if err != nil {
		return nil, err
	}
	resolver, err := NewFileResultURLResolver(cfg.ResultsRoot)
	if err != nil {
		_ = repo.Close()
		return nil, err
	}

	devAuth := NewDevAuthProvider(cfg)
	module := &Module{
		AuthProvider:      devAuth,
		JobRepository:     repo,
		ResultURLResolver: resolver,
		closers: []func() error{
			repo.Close,
		},
	}
	return module, nil
}

func (m *Module) Close() error {
	for _, closeFn := range m.closers {
		if err := closeFn(); err != nil {
			return err
		}
	}
	return nil
}
