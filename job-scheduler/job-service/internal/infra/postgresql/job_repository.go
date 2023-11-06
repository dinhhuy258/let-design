package postgresql

import (
	"context"
	"errors"
	"job-service/internal/entity"
	"job-service/internal/infra/postgresql/model"

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

func (_self *jobRepository) Create(ctx context.Context, job entity.Job) (entity.Job, error) {
	tx := _self.GetTransactionOrCreate(ctx)

	var jobModel model.Job
	jobModel = jobModel.FromEntity(job)

	err := tx.WithContext(ctx).
		Create(&jobModel).Error
	if err != nil {
		return entity.Job{}, err
	}

	return jobModel.ToEntity(), nil
}

func (_self *jobRepository) Update(ctx context.Context, job entity.Job) error {
	tx := _self.GetTransactionOrCreate(ctx)

	var jobModel model.Job
	jobModel = jobModel.FromEntity(job)

	return tx.WithContext(ctx).
		Omit("created_at").
		Save(&jobModel).Error
}

func (_self *jobRepository) FindById(ctx context.Context, id uint64) (*entity.Job, error) {
	tx := _self.GetTransactionOrCreate(ctx)

	var jobModel model.Job

	err := tx.WithContext(ctx).
		Take(&jobModel, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	jobEntity := jobModel.ToEntity()

	return &jobEntity, nil
}

func (_self *jobRepository) FindMultiByUserId(ctx context.Context, userId uint64) ([]entity.Job, error) {
	tx := _self.GetTransactionOrCreate(ctx)

	var jobModels []model.Job
	tx = tx.WithContext(ctx).Model(&jobModels)

	tx.Where("user_id = ?", userId)

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
