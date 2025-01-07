package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	challengeNameCustomChallenge = "CUSTOM_CHALLENGE"
)

func handler(ctx context.Context, event events.CognitoEventUserPoolsDefineAuthChallenge) (events.CognitoEventUserPoolsDefineAuthChallenge, error) {
	if len(event.Request.Session) > 0 &&
		hasInvalidSession(event.Request.Session) {
		// We only accept custom challenges; fail auth
		event.Response.IssueTokens = false
		event.Response.FailAuthentication = true
	} else if len(event.Request.Session) > 3 &&
		!event.Request.Session[len(event.Request.Session)-1].ChallengeResult {
		// The user provided a wrong answer 3 times; fail auth
		event.Response.IssueTokens = false
		event.Response.FailAuthentication = true
	} else if len(event.Request.Session) > 0 &&
		event.Request.Session[len(event.Request.Session)-1].ChallengeName == challengeNameCustomChallenge &&
		event.Request.Session[len(event.Request.Session)-1].ChallengeResult {
		// The user provided the right answer; succeed auth
		event.Response.IssueTokens = true
		event.Response.FailAuthentication = false
	} else {
		// The user did not provide a correct answer yet; present challenge
		event.Response.ChallengeName = challengeNameCustomChallenge
		event.Response.IssueTokens = false
		event.Response.FailAuthentication = false
	}
	return event, nil
}

func hasInvalidSession(sessions []*events.CognitoEventUserPoolsChallengeResult) bool {
	for _, s := range sessions {
		if s.ChallengeName != challengeNameCustomChallenge {
			return true
		}
	}
	return false
}

func main() {
	lambda.Start(handler)
}
