package postgres

import (
	"toy-rental-system/internal/domain/entity"
)

type SubscriptionRepository interface {
	Save(subscription *entity.Subscription) error
}

type subscriptionRepository struct {
	subscriptions map[string]*entity.Subscription
}

func NewSubscriptionRepository() SubscriptionRepository {
	return &subscriptionRepository{
		subscriptions: make(map[string]*entity.Subscription),
	}
}

func (r *subscriptionRepository) Save(subscription *entity.Subscription) error {
	r.subscriptions[subscription.ID] = subscription
	return nil
}
