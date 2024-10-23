package api

import (
	"net/http"
	controller_hello "serverless-aws-cdk/internal/controllers/hello"
	router "serverless-aws-cdk/lambdas"

	"github.com/aws/aws-lambda-go/events"
)

var HelloRoutes = map[string]router.RouteConfig{
	"/hello/": { // equivalent to /hello
		Methods: map[string]router.RouteMethodConfig{
			http.MethodGet: {
				Callback: hello_world,
			},
		},
	},
	"/hello/{goodId}": {
		Methods: map[string]router.RouteMethodConfig{
			http.MethodGet: {
				Callback: hello,
			},
		},
	},
	"/world": {
		Methods: map[string]router.RouteMethodConfig{
			http.MethodGet: {
				Callback: world,
			},
		},
	},
	"/world/tasty": {
		Methods: map[string]router.RouteMethodConfig{
			http.MethodPost: {
				Callback: world,
			},
		},
	},
}

func hello(pathParams map[string]string, addInfo router.AdditionalInfo) events.APIGatewayProxyResponse {
	goodId := pathParams["goodId"]        // should return `123` for http://localhost:4000/api/v1/test/hello/123
	bella := addInfo.QueryParams["bella"] // should return `ciao` for url http://localhost:4000/api/v1/test/hello/123?bella=ciao
	return controller_hello.Hello(goodId, bella)
}

func hello_world(pathParams map[string]string, addInfo router.AdditionalInfo) events.APIGatewayProxyResponse {
	return controller_hello.HelloWorld()
}

func world(pathParams map[string]string, addInfo router.AdditionalInfo) events.APIGatewayProxyResponse {
	return controller_hello.World()
}
