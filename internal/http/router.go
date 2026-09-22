package http

import (
	"net/http"
	"time"

	"github.com/sagemyrage/code-quality-expert-system/internal/http/handlers"
	"github.com/sagemyrage/code-quality-expert-system/internal/http/middleware"
	"github.com/sagemyrage/code-quality-expert-system/internal/service"
)

func NewRouter(
	authService *service.AuthService,
	checkService *service.CheckService,
	sessionTTL time.Duration,
	sessionCookieSecure bool,
) http.Handler {
	ah := handlers.NewAuthHandler(authService, sessionTTL, sessionCookieSecure)
	ch := handlers.NewCheckHandler(checkService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.Home)
	mux.HandleFunc("GET /health", handlers.Health)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	mux.HandleFunc("GET /login", ah.LoginPage)
	mux.HandleFunc("POST /login", ah.Login)
	mux.HandleFunc("POST /logout", ah.Logout)
	mux.HandleFunc("GET /register", ah.RegisterPage)
	mux.HandleFunc("POST /register", ah.Register)

	requireAuth := middleware.RequireAuth(http.HandlerFunc(ch.Dashboard))
	mux.Handle("GET /dashboard", requireAuth)

	createCheck := middleware.RequireAuth(http.HandlerFunc(ch.Create))
	mux.Handle("POST /checks", createCheck)

	showCheck := middleware.RequireAuth(http.HandlerFunc(ch.Show))
	mux.Handle("GET /checks/{id}", showCheck)

	identifyUser := middleware.IdentifyUser(authService)
	return identifyUser(mux)
}
