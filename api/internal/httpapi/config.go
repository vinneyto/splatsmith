package httpapi

// Config defines HTTP module settings.
type Config struct {
	ListenAddr string `yaml:"listen_addr"`
}

func (c Config) withDefaults() Config {
	if c.ListenAddr == "" {
		c.ListenAddr = ":8080"
	}
	return c
}
