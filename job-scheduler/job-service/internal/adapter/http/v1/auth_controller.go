package httpv1

import (
	"job-service/internal/entity"
	"job-service/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	HandleLogin(context *gin.Context)
}

type authController struct{}

func NewAuthController(authUsecase usecase.AuthUsecase) AuthController {
	return &authController{}
}

func (*authController) HandleLogin(context *gin.Context) {
	user := entity.User{}
	if err := context.BindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, nil)
}
