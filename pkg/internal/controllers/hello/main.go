package controller_hello

import (
	"fmt"
	"net/http"
	"serverless-aws-cdk/utils"

	"github.com/aws/aws-lambda-go/events"
)

func Hello(goodId, bella string) events.APIGatewayProxyResponse {
	fmt.Println("Got goodId & bella:", goodId, bella)

	return utils.PrepareResponse(http.StatusOK, nil, map[string]interface{}{
		"message": "Hello, World!",
	})
}

func World() events.APIGatewayProxyResponse {
	return utils.PrepareResponse(http.StatusOK, nil, map[string]interface{}{
		"message": "World, Hello!",
	})
}

func HelloWorld() events.APIGatewayProxyResponse {
	return utils.PrepareResponse(http.StatusOK, nil, map[string]interface{}{
		"message": "Should print hello_world!",
	})
}
