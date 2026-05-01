package main

import (
	"net/http"
	"goober-go/internal/handlers"
	"goober-go/internal/services"
	"goober-go/internal/repository"
	_ "goober-go/internal/middleware"
	_ "goober-go/internal/utils"
)

func registerInternalRoutes(mux *http.ServeMux) {
	handlers.NewOrderHandler("/api/order").RegisterRoutes(mux)
	handlers.NewProductHandler("/api/product").RegisterRoutes(mux)
	handlers.NewUserHandler("/api/user").RegisterRoutes(mux)
	_ = services.NewOrderService()
	_ = services.NewProductService()
	_ = services.NewUserService()
	_ = repository.NewOrderRepository()
	_ = repository.NewProductRepository()
	_ = repository.NewUserRepository()
}
