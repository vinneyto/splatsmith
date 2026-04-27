package aws

import "github.com/vinneyto/ariadne/api/internal/core"

type Module struct {
	AuthProvider   core.AuthProvider
	ScanRepository core.ScanRepository
	ObjectStorage  core.ObjectStorage
	PipelineClient core.PipelineClient
	Notifier       core.Notifier
}

func NewModule(_ Config) (*Module, error) {
	return &Module{
		AuthProvider:   &authProviderStub{},
		ScanRepository: &scanRepositoryStub{},
		ObjectStorage:  &objectStorageStub{},
		PipelineClient: &pipelineClientStub{},
		Notifier:       &notifierStub{},
	}, nil
}
