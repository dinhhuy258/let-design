package usecase

import (
	"context"
	"job-service/internal/entity"
	"job-service/pkg/logger"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, user entity.User) error
}

type userUsecase struct {
	authUsecase    AuthUsecase
	userRepository entity.UserRepository
	logger         *logger.Logger
}

func NewUserUsecase(
	authUsecase AuthUsecase,
	userRepository entity.UserRepository,
	logger *logger.Logger,
) UserUsecase {
	return &userUsecase{
		authUsecase:    authUsecase,
		userRepository: userRepository,
		logger:         logger,
	}
}

func (_self *userUsecase) CreateUser(ctx context.Context, user entity.User) error {
	hashedPwd, err := _self.authUsecase.HashPassword(user.Password)
	if err != nil {
		_self.logger.Error("failed to hash password %v", err)

		return err
	}

	user.Password = hashedPwd

	_, err = _self.userRepository.Create(ctx, user)
	if err != nil {
		_self.logger.Error("failed to create user %v", err)

		return err
	}

	return nil
}
