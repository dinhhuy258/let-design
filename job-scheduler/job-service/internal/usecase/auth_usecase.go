package usecase

import (
	"context"
	"job-service/internal/entity"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 10
)

type AuthUsecase interface {
	HashPassword(password string) (string, error)
	AttemptLogin(ctx context.Context, username, password string) (*entity.User, error)
}

type authUsecase struct {
	userRepository entity.UserRepository
}

func NewAuthUsecase(userRepository entity.UserRepository) AuthUsecase {
	return &authUsecase{
		userRepository: userRepository,
	}
}

func (_self *authUsecase) AttemptLogin(ctx context.Context, username string, password string) (*entity.User, error) {
	user, err := _self.userRepository.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, entity.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, entity.ErrAuthFailed
	}

	// clear password for security reason
	user.Password = ""

	return user, nil
}

func (*authUsecase) HashPassword(password string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(hashedPwd), nil
}
