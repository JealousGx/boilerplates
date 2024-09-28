package main

import (
	"context"
	"serverless-aws-cdk/internal/api"
	router "serverless-aws-cdk/lambdas"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return router.Router(ctx, req, api.Routes), nil
}

func main() {
	lambda.Start(handler)
}
