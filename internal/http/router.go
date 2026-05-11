package http

import (
	"net/http"

	"github.com/sagemyrage/code-quality-expert-system/internal/http/handlers"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/login", handlers.LoginPage)
	mux.HandleFunc("/register", handlers.RegisterPage)
	mux.HandleFunc("/health", handlers.Health)

	return mux
}
