package postgresql

import (
	"context"
	"job-scheduler-service/internal/entity"
	"job-scheduler-service/internal/infra/postgresql/model"
	"time"

	"gorm.io/gorm"
)

type jobRepository struct {
	transactionRepositoryInterface
	database *gorm.DB
}

func NewJobRepository(database *gorm.DB) entity.JobRepository {
	return &jobRepository{
		transactionRepositoryInterface: &transactionRepository{
			database: database,
		},
		database: database,
	}
}

func (_self *jobRepository) Update(ctx context.Context, job entity.Job) error {
	tx := _self.GetTransactionOrCreate(ctx)

	var jobModel model.Job
	jobModel = jobModel.FromEntity(job)

	return tx.WithContext(ctx).
		Omit("created_at").
		Save(&jobModel).Error
}

func (_self *jobRepository) FindMultiAvailableJobs(
	ctx context.Context,
	status string,
	executeAt time.Time,
	shardIds []uint64,
) ([]entity.Job, error) {
	return nil, nil
}
