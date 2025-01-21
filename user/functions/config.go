package functions

type Config struct {
	MailerNoReply string `env:"MAILER_NO_REPLY"`
	AuthTable     string `env:"AUTH_TABLE"`
	BasicApiUrl   string `env:"BASIC_API_URL"`
	UserPoolId    string `env:"USER_POOL_ID"`
}
