package httpv1

import (
	"job-service/pkg/httpserver"

	jwt "github.com/appleboy/gin-jwt/v2"
)

const (
	// basePath is the base path for all routes
	basePath = "/api/v1/"
)

func SetRoutes(
	server httpserver.Interface,
	authMiddleware *jwt.GinJWTMiddleware,
	userController UserController,
	jobController JobController,
) {
	router := server.GetRouter()

	apiV1Group := router.Group(basePath)
	// users
	apiV1Group.POST("users", userController.CreateUser)
	// auth
	apiV1Group.POST("auth/login", authMiddleware.LoginHandler)
	apiV1Group.POST("auth/refresh", authMiddleware.RefreshHandler)

	apiV1Group.Use(authMiddleware.MiddlewareFunc())
	{
		apiV1Group.GET("me/jobs", jobController.GetJobs)
		apiV1Group.POST("me/jobs", jobController.CreateJob)
		apiV1Group.POST("me/jobs/:job_id", jobController.CancelJob)
		apiV1Group.GET("me/jobs/:job_id", jobController.GetJob)
	}
}
