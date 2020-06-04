package mailjob

type Config struct {
	SmtpServer   string
	SmtpAddress  string
	SmtpPassword string
	SmtpPort     string
}

//smtp config ...
func NewConfig() *Config {
	return &Config{
		SmtpServer:   "smtp.gmail.com",
		SmtpAddress:  "gosmpt@gmail.com",
		SmtpPassword: "go8RmptS3",
		SmtpPort:     "587",
	}
}
