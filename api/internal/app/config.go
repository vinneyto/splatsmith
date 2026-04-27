package app

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/vinneyto/ariadne/api/internal/aws"
	"github.com/vinneyto/ariadne/api/internal/httpapi"
	"github.com/vinneyto/ariadne/api/internal/standalone"
)

type Mode string

const (
	ModeStandalone Mode = "standalone"
	ModeAWS        Mode = "aws"
)

type Config struct {
	Mode       Mode              `yaml:"mode"`
	API        httpapi.Config    `yaml:"api"`
	Standalone standalone.Config `yaml:"standalone"`
	AWS        aws.Config        `yaml:"aws"`
}

func LoadConfig(path string) (Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return Config{}, err
	}
	if cfg.Mode == "" {
		cfg.Mode = ModeStandalone
	}
	if cfg.Mode != ModeStandalone && cfg.Mode != ModeAWS {
		return Config{}, fmt.Errorf("unsupported mode: %q", cfg.Mode)
	}
	return cfg, nil
}
