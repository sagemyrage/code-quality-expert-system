package handlers

import (
	"net/http"
	"time"
)

const sessionCookieName = "session_id"

func setSessionCookie(w http.ResponseWriter, sessionID string, sessionTTL time.Duration, secure bool) {
	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessionTTL.Seconds()),
		Secure:   secure,
	}
	http.SetCookie(w, cookie)
}

func clearSessionCookie(w http.ResponseWriter, secure bool) {
	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Secure:   secure,
	}
	http.SetCookie(w, cookie)
}
