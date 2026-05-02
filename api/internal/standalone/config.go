package standalone

type Config struct {
	SQLitePath          string `yaml:"sqlite_path"`
	ResultsRoot         string `yaml:"results_root"`
	ResultURLTTLSeconds int    `yaml:"result_url_ttl_seconds"`
}
