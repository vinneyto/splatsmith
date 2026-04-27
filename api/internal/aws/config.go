package aws

type Config struct {
	Region       string `yaml:"region"`
	UserPoolID   string `yaml:"user_pool_id"`
	JobsTable    string `yaml:"jobs_table"`
	ResultBucket string `yaml:"result_bucket"`
}
