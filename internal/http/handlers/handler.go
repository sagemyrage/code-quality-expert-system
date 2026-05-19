package handlers

import (
	"time"

	"github.com/sagemyrage/code-quality-expert-system/internal/service"
)

type AuthHandler struct {
	authService         *service.AuthService
	sessionTTL          time.Duration
	sessionCookieSecure bool
}

func NewAuthHandler(authService *service.AuthService, sessionTTL time.Duration, sessionCookieSecure bool) *AuthHandler {
	return &AuthHandler{
		authService:         authService,
		sessionTTL:          sessionTTL,
		sessionCookieSecure: sessionCookieSecure,
	}
}

type CheckHandler struct {
	checkService *service.CheckService
}

func NewCheckHandler(checkService *service.CheckService) *CheckHandler {
	return &CheckHandler{
		checkService: checkService,
	}
}
