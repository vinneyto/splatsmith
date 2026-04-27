package standalone

type Config struct {
	SQLitePath          string `yaml:"sqlite_path"`
	ResultsRoot         string `yaml:"results_root"`
	ResultURLTTLSeconds int    `yaml:"result_url_ttl_seconds"`
	DevToken            string `yaml:"dev_token"`
	DevUserID           string `yaml:"dev_user_id"`
	DevUserEmail        string `yaml:"dev_user_email"`
}
