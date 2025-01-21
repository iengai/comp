package usecase

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type (
	SendOTPViaEmail interface {
		Execute(ctx context.Context, cmd *SendOTPViaEmailCommand) error
	}

	sendOTPViaEmail struct {
		mailClient    *ses.Client
		mailerNoReply string
	}

	SendOTPViaEmailCommand struct {
		emailAddress string
		code         string
	}
)

func (s sendOTPViaEmail) Execute(ctx context.Context, cmd *SendOTPViaEmailCommand) error {
	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP code is: %s", cmd.code)

	input := &ses.SendEmailInput{
		Source: aws.String(s.mailerNoReply),
		Destination: &types.Destination{
			ToAddresses: []string{cmd.emailAddress},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(subject),
			},
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(body),
				},
			},
		},
	}
	_, err := s.mailClient.SendEmail(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

func NewSendOTPViaEmail(mailClient *ses.Client, mailerNoReply string) SendOTPViaEmail {
	return &sendOTPViaEmail{
		mailClient:    mailClient,
		mailerNoReply: mailerNoReply,
	}
}

func NewSendOTPViaEmailCommand(code string, emailAddress string) *SendOTPViaEmailCommand {
	return &SendOTPViaEmailCommand{
		code:         code,
		emailAddress: emailAddress,
	}
}
