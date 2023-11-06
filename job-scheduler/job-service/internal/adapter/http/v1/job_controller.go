package httpv1

import (
	"job-service/internal/entity"
	"job-service/internal/usecase"
	"net/http"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type JobController interface {
	CreateJob(context *gin.Context)
	CancelJob(context *gin.Context)
}

type jobController struct {
	jobUsecase usecase.JobUsecase
}

func NewJobController(jobUsecase usecase.JobUsecase) JobController {
	return &jobController{
		jobUsecase: jobUsecase,
	}
}

func (_self *jobController) CreateJob(c *gin.Context) {
	job := entity.Job{}
	if err := c.BindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"success": false,
		})

		return
	}

	user, _ := c.Get(jwt.IdentityKey)
	userId := user.(*entity.User).Id
	job.UserId = userId

	job, err := _self.jobUsecase.CreateJob(c, job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"success": false,
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    job,
		"success": true,
	})
}

func (_self *jobController) CancelJob(c *gin.Context) {
	user, _ := c.Get(jwt.IdentityKey)
	userId := user.(*entity.User).Id

	jobIdParam := c.Param("job_id")

	jobId, err := strconv.ParseUint(jobIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"success": false,
		})
	}

	err = _self.jobUsecase.CancelJob(c, userId, jobId)
	if err != nil {
		// TODO: Handle error properly
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"success": false,
		})

		return
	}

	c.Status(http.StatusNoContent)
}
