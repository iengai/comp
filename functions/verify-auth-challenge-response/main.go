package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handler processes the VerifyAuthChallengeResponse trigger
func handler(ctx context.Context, event events.CognitoEventUserPoolsVerifyAuthChallenge) (events.CognitoEventUserPoolsVerifyAuthChallenge, error) {
	log.Printf("Received event: %+v\n", event)

	// Retrieve the expected answer from privateChallengeParameters
	expectedAnswer, exists := event.Request.PrivateChallengeParameters["secretLoginCode"]
	if !exists {
		log.Println("Error: 'secretLoginCode' not found in privateChallengeParameters")
		event.Response.AnswerCorrect = false
		return event, nil
	}

	// Compare the challenge answer provided by the user with the expected answer
	if event.Request.ChallengeAnswer == expectedAnswer {
		event.Response.AnswerCorrect = true
		log.Println("Challenge answer is correct.")
	} else {
		event.Response.AnswerCorrect = false
		log.Println("Challenge answer is incorrect.")
	}

	return event, nil
}

func main() {
	lambda.Start(handler)
}
