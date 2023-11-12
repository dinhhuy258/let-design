package httpv1

import (
	"job-service/internal/entity"
	"job-service/internal/usecase"
	"job-service/pkg/httpserver"
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

func (_self *userController) CreateUser(c *gin.Context) {
	user := entity.User{}
	if err := c.BindJSON(&user); err != nil {
		httpserver.ErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	err := _self.userUsecase.CreateUser(c, user)
	if err != nil {
		httpserver.ErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	httpserver.StatusResponse(c, http.StatusCreated)
}
