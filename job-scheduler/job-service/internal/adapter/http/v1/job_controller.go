package httpv1

import (
	"job-service/internal/entity"
	"job-service/internal/usecase"
	"job-service/pkg/httpserver"
	"net/http"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type JobController interface {
	GetJobs(c *gin.Context)
	CreateJob(c *gin.Context)
	CancelJob(c *gin.Context)
	GetJob(c *gin.Context)
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
		_ = c.Error(entity.ErrBadRequest)

		return
	}

	user, _ := c.Get(jwt.IdentityKey)
	userId := user.(*entity.User).Id
	job.UserId = userId

	job, err := _self.jobUsecase.CreateJob(c, job)
	if err != nil {
		_ = c.Error(err)

		return
	}

	httpserver.SuccessResponse(c, job)
}

func (_self *jobController) CancelJob(c *gin.Context) {
	user, _ := c.Get(jwt.IdentityKey)
	userId := user.(*entity.User).Id

	jobIdParam := c.Param("job_id")

	jobId, err := strconv.ParseUint(jobIdParam, 10, 64)
	if err != nil {
		_ = c.Error(entity.ErrBadRequest)

		return
	}

	err = _self.jobUsecase.CancelJob(c, userId, jobId)
	if err != nil {
		_ = c.Error(err)

		return
	}

	httpserver.StatusResponse(c, http.StatusNoContent)
}

func (_self *jobController) GetJobs(c *gin.Context) {
	user, _ := c.Get(jwt.IdentityKey)
	userId := user.(*entity.User).Id

	jobs, err := _self.jobUsecase.GetJobs(c, userId)
	if err != nil {
		_ = c.Error(err)

		return
	}

	httpserver.SuccessResponse(c, jobs)
}

func (_self *jobController) GetJob(c *gin.Context) {
	user, _ := c.Get(jwt.IdentityKey)
	userId := user.(*entity.User).Id

	jobIdParam := c.Param("job_id")

	jobId, err := strconv.ParseUint(jobIdParam, 10, 64)
	if err != nil {
		_ = c.Error(entity.ErrBadRequest)

		return
	}

	job, err := _self.jobUsecase.GetJob(c, userId, jobId)
	if err != nil {
		_ = c.Error(err)

		return
	}

	httpserver.SuccessResponse(c, job)
}
