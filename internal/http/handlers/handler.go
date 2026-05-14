package handlers

import "github.com/sagemyrage/code-quality-expert-system/internal/service"

type Handler struct {
	authService *service.AuthService
}

func NewHandler(authService *service.AuthService) *Handler {
	return &Handler{authService: authService}
}
