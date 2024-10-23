package api

import (
	"net/http"
	controller_users "serverless-aws-cdk/internal/controllers/users"
	router "serverless-aws-cdk/lambdas"
	"serverless-aws-cdk/utils"

	"github.com/aws/aws-lambda-go/events"
)

var UserRoutes = map[string]router.RouteConfig{
	"/user/{userId}": {
		Methods: map[string]router.RouteMethodConfig{
			http.MethodGet: {
				Callback: getUser,
			},
		},
	},
	"/user": {
		Methods: map[string]router.RouteMethodConfig{
			http.MethodPost: {
				Callback: createUser,
			},
		},
	},
	"/all": {
		Methods: map[string]router.RouteMethodConfig{
			http.MethodGet: {
				Callback: getAllUsers,
			},
		},
	},
}

func getUser(pathParams map[string]string, addInfo router.AdditionalInfo) events.APIGatewayProxyResponse {
	userId := pathParams["userId"]
	user, err := controller_users.GetUser(userId)

	if err != nil {
		return utils.PrepareResponse(http.StatusBadRequest, nil, map[string]interface{}{
			"message": err,
		})
	}

	return utils.PrepareResponse(http.StatusOK, nil, map[string]interface{}{
		"user": user,
	})
}

func createUser(pathParams map[string]string, addInfo router.AdditionalInfo) events.APIGatewayProxyResponse {
	userInfo := addInfo.Body
	err := controller_users.CreateUser(userInfo["name"].(string), userInfo["email"].(string), userInfo["password"].(string))

	if err != nil {
		return utils.PrepareResponse(http.StatusBadRequest, nil, map[string]interface{}{
			"message": err,
		})
	}

	return utils.PrepareResponse(http.StatusOK, nil, utils.Responses[201])
}

func getAllUsers(pathParams map[string]string, addInfo router.AdditionalInfo) events.APIGatewayProxyResponse {
	users, err := controller_users.GetAllUsers()

	if err != nil {
		return utils.PrepareResponse(http.StatusBadRequest, nil, map[string]interface{}{
			"message": err,
		})
	}

	return utils.PrepareResponse(http.StatusOK, nil, map[string]interface{}{
		"users": users,
	})
}
