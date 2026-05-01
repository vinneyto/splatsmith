package standalone

import "github.com/vinneyto/splatmaker/api/internal/core"

type Module struct {
	AuthProvider                       core.AuthProvider
	LoginProvider                      core.LoginProvider
	ReconstructionJobRepository        core.ReconstructionJobRepository
	ReconstructionSubmissionDispatcher core.ReconstructionSubmissionDispatcher
	ReconstructionResultURLResolver    core.ReconstructionResultURLResolver

	closers []func() error
}

func NewModule(cfg Config) (*Module, error) {
	repo, err := NewSQLiteReconstructionJobRepository(cfg.SQLitePath)
	if err != nil {
		return nil, err
	}
	resolver, err := NewFileResultURLResolver(cfg.ResultsRoot)
	if err != nil {
		_ = repo.Close()
		return nil, err
	}
	dispatcher := NewSimulatedReconstructionSubmissionDispatcher(repo)

	devAuth := NewDevAuthProvider(cfg)
	module := &Module{
		AuthProvider:                       devAuth,
		LoginProvider:                      devAuth,
		ReconstructionJobRepository:        repo,
		ReconstructionSubmissionDispatcher: dispatcher,
		ReconstructionResultURLResolver:    resolver,
		closers: []func() error{
			dispatcher.Close,
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
