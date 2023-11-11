package usecase

import (
	"job-scheduler-service/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type fairSchedulerUsecaseTestSuite struct {
	suite.Suite
	usecase fairSchedulerUsecase
}

func (suite *fairSchedulerUsecaseTestSuite) SetupSuite() {
	suite.usecase = fairSchedulerUsecase{}
}

func (suite *fairSchedulerUsecaseTestSuite) TearDownSuite() {
}

func (suite *fairSchedulerUsecaseTestSuite) TearDownTest() {
}

func (suite *fairSchedulerUsecaseTestSuite) TestMultipleUsersEqualWeights() {
	now := time.Now()
	user1Id := uint64(1)
	user2Id := uint64(2)
	user1Weight := float32(2)
	user2Weight := float32(2)

	jobs := suite.usecase.getScheduledJobs([]entity.Job{
		{
			Id:        1,
			ExecuteAt: now.Add(-2 * time.Second),
			User: &entity.User{
				Id:        user1Id,
				JobWeight: user1Weight,
			},
		},
		{
			Id:        2,
			ExecuteAt: now.Add(-1 * time.Second),
			User: &entity.User{
				Id:        user2Id,
				JobWeight: user2Weight,
			},
		},
		{
			Id:        3,
			ExecuteAt: now,
			User: &entity.User{
				Id:        user1Id,
				JobWeight: user1Weight,
			},
		},
		{
			Id:        4,
			ExecuteAt: now.Add(1 * time.Second),
			User: &entity.User{
				Id:        user2Id,
				JobWeight: user2Weight,
			},
		},
		{
			Id:        5,
			ExecuteAt: now.Add(2 * time.Second),
			User: &entity.User{
				Id:        user1Id,
				JobWeight: user1Weight,
			},
		},
	})

	jobIds := getJobIds(jobs)

	suite.Equal([]uint64{1, 2, 3, 4, 5}, jobIds)
}

func (suite *fairSchedulerUsecaseTestSuite) TestFairShare() {
	now := time.Now()
	user1 := &entity.User{
		Id:        1,
		JobWeight: 1,
	}
	user2 := &entity.User{
		Id:        2,
		JobWeight: 1,
	}

	jobs := suite.usecase.getScheduledJobs([]entity.Job{
		{
			Id:        1,
			ExecuteAt: now.Add(-1 * time.Second),
			User:      user1,
		},
		{
			Id:        2,
			ExecuteAt: now.Add(-2 * time.Second),
			User:      user1,
		},
		{
			Id:        3,
			ExecuteAt: now,
			User:      user1,
		},
		{
			Id:        4,
			ExecuteAt: now.Add(1 * time.Second),
			User:      user2,
		},
		{
			Id:        5,
			ExecuteAt: now.Add(2 * time.Second),
			User:      user2,
		},
	})

	jobIds := getJobIds(jobs)

	suite.Equal([]uint64{2, 4, 1, 5, 3}, jobIds)
}

func getJobIds(jobs []entity.Job) []uint64 {
	var jobIds []uint64
	for _, job := range jobs {
		jobIds = append(jobIds, job.Id)
	}

	return jobIds
}

func TestFairSchedulerUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(fairSchedulerUsecaseTestSuite))
}
