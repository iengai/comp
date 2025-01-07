package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	// add your custom logic for pre-signup
	return event, nil
}

func main() {
	lambda.Start(handler)
}
