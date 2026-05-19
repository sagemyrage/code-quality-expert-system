package handlers

import (
	"errors"
	"net/http"

	"github.com/sagemyrage/code-quality-expert-system/internal/http/requestcontext"
	"github.com/sagemyrage/code-quality-expert-system/internal/service"
)

func (h *CheckHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := requestcontext.UserID(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sourceCode := r.FormValue("code")
	_, err := h.checkService.Create(r.Context(), userID, sourceCode)
	if err != nil {
		var validationError *service.ValidationError
		if errors.As(err, &validationError) {
			http.Error(w, validationError.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
