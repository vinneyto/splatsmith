package app

import (
	"fmt"

	"github.com/vinneyto/splatmaker/api/internal/aws"
	"github.com/vinneyto/splatmaker/api/internal/core/services"
	"github.com/vinneyto/splatmaker/api/internal/standalone"
)

type Runtime struct {
	Mode         Mode
	AuthService  *services.AuthService
	JobViewer    *services.JobViewerService
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
			AuthService:  services.NewAuthService(module.AuthProvider),
			JobViewer:    services.NewJobViewerService(module.JobRepository, module.ResultURLResolver),
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
			AuthService:  services.NewAuthService(module.AuthProvider),
			JobViewer:    services.NewJobViewerService(module.JobRepository, module.ResultURLResolver),
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
