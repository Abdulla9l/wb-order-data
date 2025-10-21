package server

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"wb-order-data/internal/service"
)

type Handler struct {
	Service *service.OrderService
}

func NewHandler(s *service.OrderService) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	order, _, err := h.Service.GetOrder(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
