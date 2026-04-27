package httpapi

// Config defines HTTP module settings.
type Config struct {
	ListenAddr  string `yaml:"listen_addr"`
	OpenAPIPath string `yaml:"openapi_path"`
}

func (c Config) withDefaults() Config {
	if c.ListenAddr == "" {
		c.ListenAddr = ":8080"
	}
	if c.OpenAPIPath == "" {
		c.OpenAPIPath = "openapi.yaml"
	}
	return c
}
