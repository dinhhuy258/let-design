package httpv1

import (
	"job-service/internal/entity"
	"job-service/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	CreateUser(context *gin.Context)
}

type userController struct {
	userUsecase usecase.UserUsecase
}

func NewUserController(userUsecase usecase.UserUsecase) UserController {
	return &userController{
		userUsecase: userUsecase,
	}
}

func (_self *userController) CreateUser(context *gin.Context) {
	user := entity.User{}
	if err := context.BindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	err := _self.userUsecase.CreateUser(context, user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	context.JSON(http.StatusCreated, nil)
}
