package router

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"serverless-aws-cdk/utils"

	"github.com/aws/aws-lambda-go/events"
)

type AdditionalInfo struct {
	QueryParams map[string]string
	Body        map[string]interface{}
}

type RouteMethodConfig struct {
	Callback     func(pathParams map[string]string, addInfo AdditionalInfo) events.APIGatewayProxyResponse
	Authenticate bool // You can add more fields as needed
}

type RouteConfig struct {
	Methods map[string]RouteMethodConfig
}

func Router(ctx context.Context, req events.APIGatewayProxyRequest, routes map[string]RouteConfig) events.APIGatewayProxyResponse {
	pathParams := req.PathParameters
	addInfo := AdditionalInfo{
		QueryParams: req.QueryStringParameters,
	}
	pathname := strings.TrimSuffix(req.Path, "/") // remove trailing slash
	parts := strings.Split(pathname, "/")
	if len(parts) > 4 {
		pathname = "/" + strings.Join(parts[4:], "/")
	} else {
		pathname = "/" + strings.Join(parts[1:], "/")
	}

	httpMethod := req.HTTPMethod

	json.Unmarshal([]byte(req.Body), &addInfo.Body)

	var requestedRoute RouteConfig
	var routeExists bool

	for path, route := range routes {
		path = strings.TrimSuffix(path, "/") // remove trailing slash

		regexPattern := "^" + regexp.MustCompile(`{([^}]*)}`).ReplaceAllString(path, `(?P<$1>[^/]*)`) + "$"
		regex := regexp.MustCompile(regexPattern)
		matches := regex.FindStringSubmatch(pathname)

		if matches != nil {
			requestedRoute = route
			routeExists = true

			// need to extract path params because they are in different format by default
			// Example: http://localhost:4000/api/v1/test/hello/123, where `123` is `goodId` (/hello/{goodId}),
			// `req.PathParameters` will return proxy:hello/123
			for i, name := range regex.SubexpNames() {
				if i != 0 && name != "" {
					pathParams[name] = matches[i]
				}
			}
			break
		}

	}

	if !routeExists {
		return utils.PrepareResponse(http.StatusNotFound, nil, utils.Responses[404])
	}

	methodConfig, methodExists := requestedRoute.Methods[httpMethod]
	if !methodExists {
		return utils.PrepareResponse(http.StatusMethodNotAllowed, nil, utils.Responses[405])
	}

	return methodConfig.Callback(pathParams, addInfo)
}
