package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/caarlos0/env/v11"
	"github.com/iengai/comp/user/functions"
)

var h handler

type handler struct {
	url           string
	mailerNoReply string
	emailClient   *ses.Client
}

func init() {
	c, err := env.ParseAs[functions.Config]()
	if err != nil {
		log.Fatalf("unable to load env config: %v", err)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}
	h = handler{
		url:           c.BasicApiUrl,
		mailerNoReply: c.MailerNoReply,
		emailClient:   ses.NewFromConfig(cfg),
	}
}

func (h handler) Handle(ctx context.Context, event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	subject := "Verify your email"
	body := fmt.Sprintf("Click here to complete sign up: %s?email=%s", h.url, event.Request.UserAttributes["email"])
	input := &ses.SendEmailInput{
		Source: aws.String(h.mailerNoReply),
		Destination: &types.Destination{
			ToAddresses: []string{event.Request.UserAttributes["email"]},
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
	_, err := h.emailClient.SendEmail(ctx, input)
	if err != nil {
		return event, err
	}
	return event, nil
}

func main() {
	lambda.Start(h.Handle)
}
