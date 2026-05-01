package main

import (
	"net/http"
	"goober-go/internal/handlers"
	"goober-go/internal/services"
)

func registerInternalRoutes(mux *http.ServeMux) {
	handlers.NewUserHandler("/api/user").RegisterRoutes(mux)
	_ = services.NewUserService()
}
