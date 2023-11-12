package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"success": true,
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{
		"message": message,
		"success": false,
	})
}

func StatusResponse(c *gin.Context, code int) {
	c.Status(code)
}
