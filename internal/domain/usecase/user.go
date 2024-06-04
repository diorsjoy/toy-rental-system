package usecase

import (
	"errors"
	"toy-rental-system/internal/domain/entity"
	"toy-rental-system/internal/repository"
)

type UserUsecase interface {
	Register(user *entity.User) error
	Login(username, password string) (string, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: repo,
	}
}

func (u *userUsecase) Register(user *entity.User) error {
	// Perform business logic, such as validation, before saving
	if user.Username == "" || user.Password == "" {
		return errors.New("username and password are required")
	}
	return u.userRepo.Save(user)
}

func (u *userUsecase) Login(username, password string) (string, error) {
	user, err := u.userRepo.FindByUsername(username)
	if err != nil || user.Password != password {
		return "", errors.New("invalid credentials")
	}
	// Mock token generation for simplicity
	token := "mock-token"
	return token, nil
}
