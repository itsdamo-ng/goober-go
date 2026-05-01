package main

import (
	"net/http"
	"goober-go/internal/handlers"
	"goober-go/internal/services"
	"goober-go/internal/repository"
)

func registerInternalRoutes(mux *http.ServeMux) {
	handlers.NewCampaignHandler("/api/campaign").RegisterRoutes(mux)
	handlers.NewCommentHandler("/api/comment").RegisterRoutes(mux)
	handlers.NewInvoiceHandler("/api/invoice").RegisterRoutes(mux)
	handlers.NewNotificationHandler("/api/notification").RegisterRoutes(mux)
	handlers.NewOrderHandler("/api/order").RegisterRoutes(mux)
	handlers.NewPaymentHandler("/api/payment").RegisterRoutes(mux)
	handlers.NewProductHandler("/api/product").RegisterRoutes(mux)
	handlers.NewReviewHandler("/api/review").RegisterRoutes(mux)
	handlers.NewShipmentHandler("/api/shipment").RegisterRoutes(mux)
	handlers.NewSubscriptionHandler("/api/subscription").RegisterRoutes(mux)
	handlers.NewTicketHandler("/api/ticket").RegisterRoutes(mux)
	handlers.NewUserHandler("/api/user").RegisterRoutes(mux)

	_ = services.NewCampaignService()
	_ = services.NewCommentService()
	_ = services.NewInvoiceService()
	_ = services.NewNotificationService()
	_ = services.NewOrderService()
	_ = services.NewPaymentService()
	_ = services.NewProductService()
	_ = services.NewReviewService()
	_ = services.NewShipmentService()
	_ = services.NewSubscriptionService()
	_ = services.NewTicketService()
	_ = services.NewUserService()

	_ = repository.NewCampaignRepository()
	_ = repository.NewCommentRepository()
	_ = repository.NewInvoiceRepository()
	_ = repository.NewNotificationRepository()
	_ = repository.NewOrderRepository()
	_ = repository.NewPaymentRepository()
	_ = repository.NewProductRepository()
	_ = repository.NewReviewRepository()
	_ = repository.NewShipmentRepository()
	_ = repository.NewSubscriptionRepository()
	_ = repository.NewTicketRepository()
	_ = repository.NewUserRepository()
}
