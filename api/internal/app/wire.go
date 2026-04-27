package app

import (
	"fmt"

	"github.com/vinneyto/ariadne/api/internal/aws"
	"github.com/vinneyto/ariadne/api/internal/core"
	"github.com/vinneyto/ariadne/api/internal/standalone"
)

type Runtime struct {
	Mode        Mode
	AuthService *core.AuthService
	Close       func() error
}

func BuildRuntime(cfg Config) (*Runtime, error) {
	switch cfg.Mode {
	case ModeStandalone:
		module, err := standalone.NewModule(cfg.Standalone)
		if err != nil {
			return nil, err
		}
		return &Runtime{
			Mode:        cfg.Mode,
			AuthService: core.NewAuthService(module.AuthProvider),
			Close:       module.Close,
		}, nil
	case ModeAWS:
		module, err := aws.NewModule(cfg.AWS)
		if err != nil {
			return nil, err
		}
		return &Runtime{
			Mode:        cfg.Mode,
			AuthService: core.NewAuthService(module.AuthProvider),
			Close: func() error {
				_ = module
				return nil
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported mode: %q", cfg.Mode)
	}
}
