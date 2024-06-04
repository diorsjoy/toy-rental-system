package repository

import "toy-rental-system/internal/domain/entity"

type UserRepository interface {
	Save(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
}
