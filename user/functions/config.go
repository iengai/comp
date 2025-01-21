package functions

type Config struct {
	AuthTable     string `env:"AUTH_TABLE,required"`
	MailerNoReply string `env:"MAILER_NO_REPLY,required"`
	BasicApiUrl   string `env:"BASIC_API_URL"`
	UserPoolId    string `env:"USER_POOL_ID"`
}
