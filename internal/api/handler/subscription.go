package handler

import (
	"encoding/json"
	"fmt"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"log"
	"net/http"
	"toy-rental-system/internal/config"
	"toy-rental-system/internal/domain/entity"
	"toy-rental-system/internal/repository/postgres"
	_ "toy-rental-system/internal/service"
)

//type SubscriptionHandler struct {
//	SubscriptionService service.SubscriptionService
//}
//
//func (h *SubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
//	var sub entity.Subscription
//	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	if err := h.SubscriptionService.Subscribe(&sub); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(sub)
//}

func CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var subscription entity.Subscription
	_ = json.NewDecoder(r.Body).Decode(&subscription)

	// Process payment with Stripe
	paymentErr := ProcessPayment(subscription)
	if paymentErr != nil {
		http.Error(w, paymentErr.Error(), http.StatusInternalServerError)
		return
	}

	// Save to DB
	dbErr := postgres.SaveSubscription(subscription)
	if dbErr != nil {
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscription)
}

func ProcessPayment(subscription entity.Subscription) error {
	env, envErr := config.LoadConfig("toy-rental-system/config/app.env")
	if envErr != nil {
		log.Fatal("cannot load config:", envErr)
	}

	stripe.Key = env.StripeSecret

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(subscription.Price),
		Currency: stripe.String(subscription.Currency),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
	}

	_, err := paymentintent.New(params)
	if err != nil {
		return fmt.Errorf("failed to create payment intent: %v", err)
	}
	return nil
}
