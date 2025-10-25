package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wb-order-data/internal/models"

	"github.com/gorilla/mux"
)

type mockService struct{}

func (m *mockService) GetOrder(orderUID string) (*models.Order, error) {
	if orderUID == "test-123" {
		return &models.Order{OrderUID: "test-123", TrackNumber: "TRACK-123"}, nil
	}
	return nil, nil
}

func (m *mockService) ProcessOrderFromMessage(data []byte) error {
	return nil
}

func (m *mockService) RestoreCacheFromDB() error {
	return nil
}

func (m *mockService) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{"size": 1}
}

func TestGetOrderHandler(t *testing.T) {
	service := &mockService{}
	server := NewHTTPServer(service)

	req := httptest.NewRequest("GET", "/order/test-123", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "test-123"})

	rr := httptest.NewRecorder()
	server.GetOrderHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}

	var order models.Order
	json.NewDecoder(rr.Body).Decode(&order)
	if order.OrderUID != "test-123" {
		t.Errorf("Expected order test-123, got %s", order.OrderUID)
	}
}

func TestHealthCheckHandler(t *testing.T) {
	service := &mockService{}
	server := NewHTTPServer(service)

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	server.HealthCheckHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
}
