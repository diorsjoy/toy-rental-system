package service

//import (
//	"github.com/stripe/stripe-go/customer"
//	"github.com/stripe/stripe-go/sub"
//	"github.com/stripe/stripe-go/v72"
//	"toy-rental-system/internal/config"
//)
//
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
