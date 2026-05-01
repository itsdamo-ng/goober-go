package main

import (
	"net/http"
	_ "goober-go/internal/models"
)

func registerInternalRoutes(mux *http.ServeMux) {
	_ = mux
}
