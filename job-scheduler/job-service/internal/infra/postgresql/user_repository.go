package postgresql

import (
	"context"
	"errors"
	"job-service/internal/entity"
	"job-service/internal/infra/postgresql/model"

	"gorm.io/gorm"
)

type userRepository struct {
	transactionRepositoryInterface
	database *gorm.DB
}

func NewUserRepository(database *gorm.DB) entity.UserRepository {
	return &userRepository{
		transactionRepositoryInterface: &transactionRepository{
			database: database,
		},
		database: database,
	}
}

func (_self *userRepository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	tx := _self.GetTransactionOrCreate(ctx)

	var userModel model.User
	userModel = userModel.FromEntity(user)

	err := tx.WithContext(ctx).
		Create(&userModel).Error
	if err != nil {
		return entity.User{}, err
	}

	return userModel.ToEntity(), nil
}

func (_self *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	tx := _self.GetTransactionOrCreate(ctx)

	var user model.User

	err := tx.WithContext(ctx).Take(&user, "username = ?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	userEntity := user.ToEntity()

	return &userEntity, nil
}
