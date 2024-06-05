package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"toy-rental-system/internal/domain/entity"
)

type SubscriptionRepository interface {
	Save(subscription *entity.Subscription) error
}

type subscriptionRepository struct {
	DB *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) SubscriptionRepository {
	return &subscriptionRepository{
		DB: db,
	}
}

func (r *subscriptionRepository) Save(subscription *entity.Subscription) error {
	query := `INSERT INTO subscriptions (id, user_id, tokens, price, currency) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Exec(query, subscription.ID, subscription.UserID, subscription.Tokens, subscription.Price, subscription.Currency)
	return err
}
