package postgresql

import (
	"context"
	"job-scheduler-worker/internal/entity"
	"job-scheduler-worker/internal/infra/postgresql/model"
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
	tx := _self.GetTransactionOrCreate(ctx)

	var jobModels []model.Job
	tx = tx.WithContext(ctx).Model(&jobModels)
	tx = tx.Where("status = ?", status)
	tx = tx.Where("execute_at <= ?", executeAt)
	tx = tx.Where("shard_id IN ?", shardIds)
	tx = tx.Preload("User")

	err := tx.Find(&jobModels).Error
	if err != nil {
		return nil, err
	}

	var jobs []entity.Job
	for _, jobModel := range jobModels {
		jobs = append(jobs, jobModel.ToEntity())
	}

	return jobs, nil
}
