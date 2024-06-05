package handler

import (
	"encoding/json"
	"net/http"
	"toy-rental-system/internal/domain/entity"
	"toy-rental-system/internal/service"
	_ "toy-rental-system/internal/service"
)

type SubscriptionHandler struct {
	SubscriptionService service.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		SubscriptionService: *subscriptionService,
	}
}

func (h *SubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var sub entity.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process payment with Stripe
	if err := h.SubscriptionService.ProcessPayment(&sub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.SubscriptionService.Subscribe(&sub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sub)
}
