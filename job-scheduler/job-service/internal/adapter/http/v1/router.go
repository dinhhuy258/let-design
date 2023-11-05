package httpv1

import (
	"job-service/pkg/httpserver"
)

const (
	// basePath is the base path for all routes
	basePath = "/api/v1/"
)

func SetRoutes(
	server httpserver.Interface,
	authController AuthController,
	userController UserController,
) {
	router := server.GetRouter()

	apiV1Group := router.Group(basePath)
	// users
	apiV1Group.POST("users", userController.CreateUser)
	// auth
	apiV1Group.POST("login", authController.HandleLogin)
}
