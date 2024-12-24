package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

var (
	cognitoClient *cognitoidentityprovider.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}
	cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)
}

func handler(ctx context.Context, event events.CognitoEventUserPoolsPostAuthentication) (events.CognitoEventUserPoolsPostAuthentication, error) {
	log.Printf("Received event: %+v\n", event)

	emailVerified, exists := event.Request.UserAttributes["email_verified"]
	if !exists || emailVerified != "true" {
		params := &cognitoidentityprovider.AdminUpdateUserAttributesInput{
			UserPoolId: aws.String(event.UserPoolID),
			Username:   aws.String(event.UserName),
			UserAttributes: []types.AttributeType{
				{
					Name:  aws.String("email_verified"),
					Value: aws.String("true"),
				},
			},
		}

		_, err := cognitoClient.AdminUpdateUserAttributes(ctx, params)
		if err != nil {
			log.Printf("Failed to update user attributes: %v", err)
			return event, err
		}
		log.Println("Successfully updated email_verified to true")
	}

	return event, nil
}

func main() {
	lambda.Start(handler)
}
