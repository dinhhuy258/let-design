package middleware

import (
	"errors"
	"job-service/internal/entity"
	"job-service/pkg/httpserver"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	for _, err := range c.Errors {
		if errors.Is(err.Err, entity.ErrUserNotFound) {
			httpserver.ErrorResponse(c, http.StatusNotFound, err.Error())

			return
		} else if errors.Is(err.Err, entity.ErrJobNotFound) {
			httpserver.ErrorResponse(c, http.StatusNotFound, err.Error())

			return
		} else if errors.Is(err.Err, entity.ErrAuthFailed) {
			httpserver.ErrorResponse(c, http.StatusUnauthorized, err.Error())

			return
		} else if errors.Is(err.Err, entity.ErrJobCannotBeCancelled) {
			httpserver.ErrorResponse(c, http.StatusBadRequest, err.Error())

			return
		} else if errors.Is(err.Err, entity.ErrBadRequest) {
			httpserver.ErrorResponse(c, http.StatusBadRequest, err.Error())

			return
		}
	}

	httpserver.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
}
