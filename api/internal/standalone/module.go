package standalone

import "github.com/vinneyto/splatmaker/api/internal/core"

type Module struct {
	AuthProvider               core.AuthProvider
	LoginProvider              core.LoginProvider
	JobRepository              core.JobRepository
	PipelineSettingsRepository core.PipelineSettingsRepository
	JobDispatcher              core.JobDispatcher
	ResultURLResolver          core.ResultURLResolver

	closers []func() error
}

func NewModule(cfg Config) (*Module, error) {
	repo, err := NewSQLiteJobRepository(cfg.SQLitePath)
	if err != nil {
		return nil, err
	}
	pipelineSettingsRepo, err := NewSQLitePipelineSettingsRepository(cfg.SQLitePath)
	if err != nil {
		_ = repo.Close()
		return nil, err
	}
	resolver, err := NewFileResultURLResolver(cfg.ResultsRoot)
	if err != nil {
		_ = pipelineSettingsRepo.Close()
		_ = repo.Close()
		return nil, err
	}
	dispatcher := NewSimulatedJobDispatcher(repo)

	devAuth := NewDevAuthProvider(cfg)
	module := &Module{
		AuthProvider:               devAuth,
		LoginProvider:              devAuth,
		JobRepository:              repo,
		PipelineSettingsRepository: pipelineSettingsRepo,
		JobDispatcher:              dispatcher,
		ResultURLResolver:          resolver,
		closers: []func() error{
			dispatcher.Close,
			pipelineSettingsRepo.Close,
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
