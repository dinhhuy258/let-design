package usecase

import (
	"job-service/internal/entity"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 10
)

type AuthUsecase interface {
	HashPassword(password string) (string, error)
	AttemptLogin(username, password string) (*entity.User, error)
}

type authUsecase struct {
	userRepository entity.UserRepository
}

func NewAuthUsecase(userRepository entity.UserRepository) AuthUsecase {
	return &authUsecase{
		userRepository: userRepository,
	}
}

func (*authUsecase) AttemptLogin(username string, password string) (*entity.User, error) {
	panic("unimplemented")
}

func (*authUsecase) HashPassword(password string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(hashedPwd), nil
}
