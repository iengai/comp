package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.CognitoEventUserPoolsDefineAuthChallenge) (events.CognitoEventUserPoolsDefineAuthChallenge, error) {
	if len(event.Request.Session) == 0 {
		event.Response.ChallengeName = "CUSTOM_CHALLENGE"
		event.Response.IssueTokens = false
		event.Response.FailAuthentication = false
	} else if event.Request.Session[len(event.Request.Session)-1].ChallengeResult {
		event.Response.IssueTokens = true
		event.Response.FailAuthentication = false
	} else {
		event.Response.IssueTokens = false
		event.Response.FailAuthentication = true
	}
	return event, nil
}

func main() {
	lambda.Start(handler)
}
