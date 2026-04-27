package aws

type Config struct {
	Region           string `yaml:"region"`
	UserPoolID       string `yaml:"user_pool_id"`
	ScansTable       string `yaml:"scans_table"`
	InputBucket      string `yaml:"input_bucket"`
	OutputBucket     string `yaml:"output_bucket"`
	StepFunctionsARN string `yaml:"step_functions_arn"`
	SESSenderEmail   string `yaml:"ses_sender_email"`
}
