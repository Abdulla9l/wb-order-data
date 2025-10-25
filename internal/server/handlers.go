package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *HTTPServer) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["id"]

	if orderUID == "" {
		http.Error(w, `{"error": "Order ID is required"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Received request for order: %s", orderUID)

	order, err := h.service.GetOrder(orderUID)
	if err != nil {
		log.Printf("Error getting order %s: %v", orderUID, err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	if order == nil {
		http.Error(w, `{"error": "Order not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully returned order: %s", orderUID)
}

func (h *HTTPServer) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "order-service",
	})
}

func (h *HTTPServer) CacheStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats := h.service.GetCacheStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
