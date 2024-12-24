package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handler is the entry point for the Lambda function
func handler(ctx context.Context, event events.CognitoEventUserPoolsCreateAuthChallenge) (events.CognitoEventUserPoolsCreateAuthChallenge, error) {
	log.Printf("Received event: %+v\n", event)

	// Check if this is the first challenge attempt or if the last challenge was successful
	if len(event.Request.Session) == 0 || event.Request.Session[len(event.Request.Session)-1].ChallengeResult {
		// Generate a 6-digit OTP
		otp := generateOTP()
		log.Printf("Generated OTP: %s\n", otp)

		// Extract the email from user attributes
		email := event.Request.UserAttributes["email"]
		if email == "" {
			return event, fmt.Errorf("email attribute is missing")
		}

		// Simulate sending the OTP to the user's email
		if err := sendOTP(email, otp); err != nil {
			log.Printf("Failed to send OTP: %v\n", err)
			return event, err
		}

		// Set the public and private challenge parameters
		event.Response.PublicChallengeParameters = map[string]string{
			"email": email,
		}
		event.Response.PrivateChallengeParameters = map[string]string{
			"otp": otp,
		}
		event.Response.ChallengeMetadata = fmt.Sprintf("CODE-%s", otp)
	} else {
		// If the last challenge was unsuccessful, no new challenge is generated
		event.Response.PublicChallengeParameters = map[string]string{}
		event.Response.PrivateChallengeParameters = map[string]string{}
		event.Response.ChallengeMetadata = ""
	}

	return event, nil
}

// generateOTP generates a random 6-digit OTP
func generateOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}

// sendOTP simulates sending an OTP to the user's email
func sendOTP(email, otp string) error {
	log.Printf("Sending OTP %s to email: %s\n", otp, email)
	// Placeholder for actual email-sending logic (e.g., AWS SES)
	return nil
}

func main() {
	lambda.Start(handler)
}
