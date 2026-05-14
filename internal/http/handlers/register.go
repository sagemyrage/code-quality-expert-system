package handlers

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/sagemyrage/code-quality-expert-system/internal/service"
)

func (h *Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles(
		"web/templates/base.html",
		"web/templates/register.html",
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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordConfirmation := r.FormValue("password_confirmation")

	_, err := h.authService.Register(r.Context(), email, password, passwordConfirmation)
	if err != nil {
		var validationError *service.ValidationError
		if errors.As(err, &validationError) {
			http.Error(w, validationError.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
