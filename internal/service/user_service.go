package service

import (
	"errors"
	"toy-rental-system/internal/domain/entity"
	"toy-rental-system/internal/repository"
)

type UserService interface {
	Register(user *entity.User) error
	Login(username, password string) (string, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		userRepository: repo,
	}
}

func (s *userService) Register(user *entity.User) error {
	return s.userRepository.Save(user)
}

func (s *userService) Login(username, password string) (string, error) {
	user, err := s.userRepository.FindByUsername(username)
	if err != nil || user.Password != password {
		return "", errors.New("invalid credentials")
	}
	// Mock token generation for simplicity
	token := "mock-token"
	return token, nil
}

