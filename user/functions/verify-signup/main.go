package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/caarlos0/env/v11"
	"github.com/iengai/comp/user/functions"
)

var h handler

type handler struct {
	cognitoClient *cognitoidentityprovider.Client
	userPoolID    string
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
		cognitoClient: cognitoidentityprovider.NewFromConfig(cfg),
		userPoolID:    c.UserPoolId,
	}
}

func (h handler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// todo use extra info to confirm signup
	email, ok := req.QueryStringParameters["email"]
	if !ok || email == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message":"email is required"}`,
		}, nil
	}
	user, err := h.cognitoClient.AdminGetUser(
		ctx,
		&cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(h.userPoolID),
			Username:   aws.String(email),
		})
	if err != nil {
		slog.Error("unable to get user", slog.String("err", err.Error()))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message":"get user failed","error":"%s"}`,
		}, nil
	}
	var emailVerified, userConfirmed bool
	for _, attr := range user.UserAttributes {
		if *attr.Name == "email_verified" {
			emailVerified = *attr.Value == "true"
			continue
		}
		if *attr.Name == "user_confirmed" {
			userConfirmed = *attr.Value == "true"
			break
		}
	}
	if userConfirmed {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       `{"message":"sign up already verified"}`,
		}, nil
	}

	if !emailVerified {
		_, err = h.cognitoClient.AdminUpdateUserAttributes(ctx,
			&cognitoidentityprovider.AdminUpdateUserAttributesInput{
				UserPoolId: aws.String(h.userPoolID),
				Username:   aws.String(email),
				UserAttributes: []types.AttributeType{
					{
						Name:  aws.String("email_verified"),
						Value: aws.String("true"),
					},
				},
			})
		if err != nil {
			slog.Error("unable to verify email", slog.String("err", err.Error()))
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       `{"message":"email verification failed"}`,
			}, nil
		}
	}

	_, err = h.cognitoClient.AdminConfirmSignUp(ctx,
		&cognitoidentityprovider.AdminConfirmSignUpInput{
			UserPoolId: aws.String(h.userPoolID),
			Username:   aws.String(email),
		})
	if err != nil {
		slog.Error("unable to confirm user", slog.String("err", err.Error()))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"message":"user confirmation failed"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message":"sign up successfully"}`,
	}, nil
}

func main() {
	lambda.Start(h.Handle)
}
