package functions

type Config struct {
	MailerNoReply string `env:"MAILER_NO_REPLY,required"`
}
