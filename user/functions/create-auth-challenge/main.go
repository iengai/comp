package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/caarlos0/env/v11"
	"github.com/iengai/comp/user/functions"
	"github.com/iengai/comp/user/internal/domain"
	"github.com/iengai/comp/user/internal/usecase"
)

var (
	sesClient     *ses.Client
	mailerNoReply string
)

var h handler

type handler struct {
	sendOTP usecase.SendOTPViaEmail
}

func init() {
	c, err := env.ParseAs[functions.Config]()
	if err != nil {
		log.Fatalf("unable to load env config: %v", err)
	}
	mailerNoReply = c.MailerNoReply
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	sesClient = ses.NewFromConfig(cfg)
	h = handler{
		sendOTP: usecase.NewSendOTPViaEmail(ses.NewFromConfig(cfg), c.MailerNoReply),
	}
}

// Handle is the entry point for the Lambda function
func (h handler) Handle(ctx context.Context, event events.CognitoEventUserPoolsCreateAuthChallenge) (events.CognitoEventUserPoolsCreateAuthChallenge, error) {
	log.Printf("Received event: %+v\n", event)

	// Extract the email from user attributes
	email := event.Request.UserAttributes["email"]
	if email == "" {
		return event, fmt.Errorf("email attribute is missing")
	}
	var otp string
	// Check if this is the first challenge attempt or if the last challenge was successful
	if len(event.Request.Session) > 0 {
		// Generate a 6-digit OTP
		otp = domain.NewCode(time.Now()).ToString()
		log.Printf("Generated OTP: %s\n", otp)

		// Simulate sending the OTP to the user's email
		if err := h.sendOTP.Execute(ctx, usecase.NewSendOTPViaEmailCommand(otp, email)); err != nil {
			log.Printf("Failed to send OTP: %v\n", err)
			return event, err
		}
	} else {
		preChallenge := event.Request.Session[len(event.Request.Session)-1]
		re := regexp.MustCompile(`CODE-(\d*)`)
		matches := re.FindStringSubmatch(preChallenge.ChallengeMetadata)
		otp = matches[1]
		if len(matches) < 2 {
			return event, errors.New("no match found for CODE-(\\d*)")
		}
	}

	event.Response.PublicChallengeParameters = map[string]string{
		"email": email,
	}
	event.Response.PrivateChallengeParameters = map[string]string{
		"otp": otp,
	}
	event.Response.ChallengeMetadata = fmt.Sprintf("CODE-%s", otp)
	return event, nil
}

func main() {
	lambda.Start(h.Handle)
}
