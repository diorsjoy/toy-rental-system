package service

import (
	"toy-rental-system/internal/domain/entity"
	"toy-rental-system/internal/repository/postgres"
)

type SubscriptionService interface {
	Subscribe(subscription *entity.Subscription) error
}

type subscriptionService struct {
	subscriptionRepository postgres.SubscriptionRepository
	userService            UserService
}

func NewSubscriptionService(repo postgres.SubscriptionRepository, us UserService) SubscriptionService {
	return &subscriptionService{
		subscriptionRepository: repo,
		userService:            us,
	}
}

func (s *subscriptionService) Subscribe(subscription *entity.Subscription) error {
	user, err := s.userService.GetUserByID(subscription.UserID)
	if err != nil {
		return err
	}

	// Add tokens based on subscription plan
	user.Tokens += subscription.Tokens
	if err := s.userService.Update(user); err != nil {
		return err
	}

	return s.subscriptionRepository.Save(subscription)
}
