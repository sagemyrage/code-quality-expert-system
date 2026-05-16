package handlers

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/sagemyrage/code-quality-expert-system/internal/service"
)

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles(
		"web/templates/base.html",
		"web/templates/login.html",
	)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	result, err := h.authService.Login(r.Context(), email, password)
	if err != nil {
		var validationError *service.ValidationError
		if errors.As(err, &validationError) {
			http.Error(w, validationError.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, result.SessionID, h.sessionTTL, h.sessionCookieSecure)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			clearSessionCookie(w, h.sessionCookieSecure)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = h.authService.Logout(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	clearSessionCookie(w, h.sessionCookieSecure)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
