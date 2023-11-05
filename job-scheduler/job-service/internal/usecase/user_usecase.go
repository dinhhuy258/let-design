package usecase

import "job-service/internal/entity"

type UserUsecase interface{}

type userUsecase struct {
	userRepository entity.UserRepository
}

func NewUserUsecase(userRepository entity.UserRepository) UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}
