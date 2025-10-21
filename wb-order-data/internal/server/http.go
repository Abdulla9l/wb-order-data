package server

import (
	"fmt"
	"log"
	"net/http"

	"wb-order-data/internal/service"

	"github.com/gorilla/mux"
)

type Server struct {
	orderService *service.OrderService
	port         string
}

func NewServer(orderService *service.OrderService, port string) *Server {
	return &Server{
		orderService: orderService,
		port:         port,
	}
}

func (s *Server) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/order/{id}", s.getOrderHandler).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	fmt.Printf("🚀 HTTP server started on :%s\n", s.port)
	log.Fatal(http.ListenAndServe(":"+s.port, router))
}
