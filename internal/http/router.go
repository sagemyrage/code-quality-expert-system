package http

import (
	"net/http"

	"github.com/sagemyrage/code-quality-expert-system/internal/http/handlers"
	"github.com/sagemyrage/code-quality-expert-system/internal/service"
)

func NewRouter(authService *service.AuthService) http.Handler {
	h := handlers.NewHandler(authService)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.Home)
	mux.HandleFunc("GET /login", h.LoginPage)
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("GET /register", h.RegisterPage)
	mux.HandleFunc("POST /register", h.Register)
	mux.HandleFunc("GET /health", handlers.Health)

	return mux
}
