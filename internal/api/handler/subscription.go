package handler

import (
	"encoding/json"
	"net/http"
	"toy-rental-system/internal/domain/entity"
	"toy-rental-system/internal/service"

	"github.com/gorilla/mux"
)

type SubscriptionHandler struct {
	subscriptionService service.SubscriptionService
}

func NewSubscriptionHandler(r *mux.Router, ss service.SubscriptionService) {
}

func (h *SubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var sub entity.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.subscriptionService.Subscribe(&sub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
