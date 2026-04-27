package standalone

type Config struct {
	SQLitePath   string `yaml:"sqlite_path"`
	StorageRoot  string `yaml:"storage_root"`
	DevToken     string `yaml:"dev_token"`
	DevUserID    string `yaml:"dev_user_id"`
	DevUserEmail string `yaml:"dev_user_email"`
}
