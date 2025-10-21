package server

import (
    "net/http"
    "wb-order-data/internal/service"
    "github.com/gorilla/mux"
    
)

func StartHTTPServer(service *service.OrderService, port string) {
    r := mux.NewRouter()
    handler := NewHandler(service)
    handler.RegisterRoutes(r)

    http.ListenAndServe(":"+port, r)
}
