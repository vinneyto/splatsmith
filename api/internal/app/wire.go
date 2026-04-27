package app

import (
	"fmt"

	"github.com/vinneyto/ariadne/api/internal/aws"
	"github.com/vinneyto/ariadne/api/internal/core"
	"github.com/vinneyto/ariadne/api/internal/standalone"
)

type Runtime struct {
	Mode         Mode
	AuthService  *core.AuthService
	JobViewer    *core.JobViewerService
	ResultURLTTL int
	Close        func() error
}

func BuildRuntime(cfg Config) (*Runtime, error) {
	switch cfg.Mode {
	case ModeStandalone:
		module, err := standalone.NewModule(cfg.Standalone)
		if err != nil {
			return nil, err
		}
		ttl := cfg.Standalone.ResultURLTTLSeconds
		if ttl <= 0 {
			ttl = 900
		}
		return &Runtime{
			Mode:         cfg.Mode,
			AuthService:  core.NewAuthService(module.AuthProvider),
			JobViewer:    core.NewJobViewerService(module.JobRepository, module.ResultURLResolver),
			ResultURLTTL: ttl,
			Close:        module.Close,
		}, nil
	case ModeAWS:
		module, err := aws.NewModule(cfg.AWS)
		if err != nil {
			return nil, err
		}
		return &Runtime{
			Mode:         cfg.Mode,
			AuthService:  core.NewAuthService(module.AuthProvider),
			JobViewer:    core.NewJobViewerService(module.JobRepository, module.ResultURLResolver),
			ResultURLTTL: 900,
			Close: func() error {
				_ = module
				return nil
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported mode: %q", cfg.Mode)
	}
}
