package database

//Config of database
type Config struct {
	DatabaseURL  string
	DatabaseName string
	UserColName  string
	LogColName   string
}

//NewConfig returns initialized database config
func NewConfig() *Config {

	return &Config{
		DatabaseURL:  "mongodb://localhost:27017/smptServer",
		DatabaseName: "smptServer",
		UserColName:  "users",
		LogColName:   "logs",
	}
}
