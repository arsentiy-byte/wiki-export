package config

type Storage interface {
	GetHost() string
	GetPort() int
	GetUser() string
	GetPassword() string
	GetDatabase() string
	GetMigrationFile() string
}
