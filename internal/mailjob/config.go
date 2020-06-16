package mailjob

//Config of smtp
type Config struct {
	SMTPServer   string
	SMTPAddress  string
	SMTPPassword string
	SMTPPort     string
}

//NewConfig returns initialized smtp config
func NewConfig() *Config {

	return &Config{
		SMTPServer:   "smtp.gmail.com",
		SMTPAddress:  "gosmpt@gmail.com",
		SMTPPassword: "go8RmptS3",
		SMTPPort:     "587",
	}

}
