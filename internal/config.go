package internal

import (
	"github.com/anastasja-hunko/smptServer/internal/database"
)

type Config struct {
	Port     string
	LogLevel string
	DbConfig *database.Config
}

func NewConfig() *Config {
	return &Config{
		Port:     ":8283",
		LogLevel: "debug",
		DbConfig: database.NewConfig(),
	}
}
