package standalone

import "github.com/vinneyto/ariadne/api/internal/core"

type Module struct {
	AuthProvider   core.AuthProvider
	ScanRepository core.ScanRepository
	ObjectStorage  core.ObjectStorage
	PipelineClient core.PipelineClient
	Notifier       core.Notifier

	closers []func() error
}

func NewModule(cfg Config) (*Module, error) {
	repo, err := NewSQLiteScanRepository(cfg.SQLitePath)
	if err != nil {
		return nil, err
	}
	storage, err := NewLocalObjectStorage(cfg.StorageRoot)
	if err != nil {
		_ = repo.Close()
		return nil, err
	}

	module := &Module{
		AuthProvider:   NewDevAuthProvider(cfg),
		ScanRepository: repo,
		ObjectStorage:  storage,
		PipelineClient: NewPipelineStub(),
		Notifier:       NewLogNotifier(),
		closers:        []func() error{repo.Close},
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
