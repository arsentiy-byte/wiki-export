package config

type Storage struct {
	Host          string `yaml:"host" env-default:"127.0.0.1"`
	Port          int    `yaml:"port" env-default:"3306"`
	User          string `yaml:"user" env-required:"true"`
	Password      string `yaml:"password" env-required:"true"`
	Database      string `yaml:"database" env-required:"true"`
	MigrationFile string `yaml:"migration"`
}

func (s *Storage) GetHost() string {
	return s.Host
}

func (s *Storage) GetPort() int {
	return s.Port
}

func (s *Storage) GetUser() string {
	return s.User
}

func (s *Storage) GetPassword() string {
	return s.Password
}

func (s *Storage) GetDatabase() string {
	return s.Database
}

func (s *Storage) GetMigrationFile() string {
	return s.MigrationFile
}
