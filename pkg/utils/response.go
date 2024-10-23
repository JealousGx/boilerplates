package utils

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

var Responses = map[int]map[string]interface{}{
	201: {
		"message": "Created",
	},
	204: {
		"message": "No Content",
	},
	400: {
		"message": "Bad Request",
	},
	401: {
		"message": "Unauthorized",
	},
	403: {
		"message": "Forbidden",
	},
	404: {
		"message": "Not Found",
	},
	405: {
		"message": "Method Not Allowed",
	},
	500: {
		"message": "Internal Server Error",
	},
	502: {
		"message": "Bad Gateway",
	},
	503: {
		"message": "Service Unavailable",
	},
	504: {
		"message": "Gateway Timeout",
	},
}

func PrepareResponse(statusCode int, headers map[string]string, body map[string]interface{}) events.APIGatewayProxyResponse {
	if headers == nil {
		headers = make(map[string]string)
	}

	if headers["Content-Type"] == "" {
		headers["Content-Type"] = "application/json"
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		bodyJson = []byte(`{"message": "Internal Server Error"}`)
		statusCode = 500
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       string(bodyJson),
	}
}
