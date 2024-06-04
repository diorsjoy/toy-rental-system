package service

import (
	"fmt"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"toy-rental-system/internal/config"
	"toy-rental-system/internal/domain/entity"
	"toy-rental-system/internal/repository/postgres"
)

type SubscriptionService struct {
	cfg              config.Config
	subscriptionRepo postgres.SubscriptionRepository
}

func NewSubscriptionService(cfg config.Config, subscriptionRepo postgres.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		cfg:              cfg,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *SubscriptionService) Subscribe(subscription *entity.Subscription) error {
	return s.subscriptionRepo.Save(subscription)
}

func (s *SubscriptionService) ProcessPayment(subscription *entity.Subscription) error {
	stripe.Key = s.cfg.StripeSecret

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(subscription.Price),
		Currency:           stripe.String(subscription.Currency),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
	}

	_, err := paymentintent.New(params)
	if err != nil {
		return fmt.Errorf("failed to create payment intent: %v", err)
	}
	return nil
}

//type SubscriptionService struct {
//	cfg *config.Config
//}
//
//func NewSubscriptionService(cfg *config.Config) *SubscriptionService {
//	return &SubscriptionService{cfg: cfg}
//}
//
//func (s *SubscriptionService) CreateCustomer(email, paymentMethodID string) (*stripe.Customer, error) {
//	stripe.Key = s.cfg.Stripe.APIKey
//	params := &stripe.CustomerParams{
//		Email: stripe.String(email),
//		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
//			DefaultPaymentMethod: stripe.String(paymentMethodID),
//		},
//	}
//	return customer.New(params)
//}
//
//func (s *SubscriptionService) CreateSubscription(customerID, priceID string) (*stripe.Subscription, error) {
//	stripe.Key = s.cfg.Stripe.APIKey
//	params := &stripe.SubscriptionParams{
//		Customer: stripe.String(customerID),
//		Items: []*stripe.SubscriptionItemsParams{
//			{
//				Price: stripe.String(priceID),
//			},
//		},
//	}
//	return sub.New(params)
//}
