package server

import (
	"log"
	"net/http"
	"wb-order-data/internal/service"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	service *service.OrderService
	router  *mux.Router
}

func NewHTTPServer(service *service.OrderService) *HTTPServer {
	server := &HTTPServer{
		service: service,
		router:  mux.NewRouter(),
	}

	server.setupRoutes()
	return server
}

func (h *HTTPServer) setupRoutes() {
	h.router.HandleFunc("/order/{id}", h.GetOrderHandler).Methods("GET")
	h.router.HandleFunc("/health", h.HealthCheckHandler).Methods("GET")
	h.router.HandleFunc("/cache/stats", h.CacheStatsHandler).Methods("GET")
	h.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))
}

func (h *HTTPServer) Start(addr string) error {
	log.Printf("Starting HTTP server on %s", addr)
	return http.ListenAndServe(addr, h.router)
}
