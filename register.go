package main

import (
	"net/http"
	"goober-go/internal/handlers"
	"goober-go/internal/services"
	"goober-go/internal/repository"
)

func registerInternalRoutes(mux *http.ServeMux) {
	handlers.NewInvoiceHandler("/api/invoice").RegisterRoutes(mux)
	handlers.NewOrderHandler("/api/order").RegisterRoutes(mux)
	handlers.NewPaymentHandler("/api/payment").RegisterRoutes(mux)
	handlers.NewProductHandler("/api/product").RegisterRoutes(mux)
	handlers.NewShipmentHandler("/api/shipment").RegisterRoutes(mux)
	handlers.NewUserHandler("/api/user").RegisterRoutes(mux)

	_ = services.NewInvoiceService()
	_ = services.NewOrderService()
	_ = services.NewPaymentService()
	_ = services.NewProductService()
	_ = services.NewShipmentService()
	_ = services.NewUserService()

	_ = repository.NewInvoiceRepository()
	_ = repository.NewOrderRepository()
	_ = repository.NewPaymentRepository()
	_ = repository.NewProductRepository()
	_ = repository.NewShipmentRepository()
	_ = repository.NewUserRepository()
}
