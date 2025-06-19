package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tahsin005/codercat-server/domain"
	"github.com/tahsin005/codercat-server/service"
)

type SubscriberHandler struct {
	service service.SubscriberService
}

func NewSubscriberHandler(service service.SubscriberService) *SubscriberHandler {
	return &SubscriberHandler{service: service}
}

func (h *SubscriberHandler) CreateSubscriber(w http.ResponseWriter, r *http.Request) {
	var subscriber domain.Subscriber
	if err := json.NewDecoder(r.Body).Decode(&subscriber); err != nil || subscriber.Email == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.CreateSubscriber(r.Context(), &subscriber); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Subscribed successfully",
	})
}

func (h *SubscriberHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/subscribe", h.CreateSubscriber).Methods("POST")
}
