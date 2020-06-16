package internal

import (
	"github.com/anastasja-hunko/smptServer/internal/database"
)

//Config of system
type Config struct {
	Port     string
	LogLevel string
	DbConfig *database.Config
}

//NewConfig returns initialized system gonfig
func NewConfig() *Config {

	return &Config{
		Port:     ":8283",
		LogLevel: "debug",
		DbConfig: database.NewConfig(),
	}
}
