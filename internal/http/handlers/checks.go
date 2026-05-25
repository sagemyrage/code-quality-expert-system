package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/sagemyrage/code-quality-expert-system/internal/domain"
	"github.com/sagemyrage/code-quality-expert-system/internal/http/requestcontext"
	"github.com/sagemyrage/code-quality-expert-system/internal/repository"
	"github.com/sagemyrage/code-quality-expert-system/internal/service"
)

type CheckPageData struct {
	Check           domain.Check
	Metrics         domain.CheckMetrics
	Recommendations []domain.CheckRecommendation
}

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

	sourceCode := r.FormValue("source_code")
	check, err := h.checkService.Create(r.Context(), userID, sourceCode)
	if err != nil {
		var validationError *service.ValidationError
		if errors.As(err, &validationError) {
			http.Error(w, validationError.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/checks/%d", check.ID), http.StatusSeeOther)
}

func (h *CheckHandler) Show(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := requestcontext.UserID(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id := r.PathValue("id")
	checkID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	details, err := h.checkService.GetDetailsByID(r.Context(), checkID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrCheckNotFound) {
			http.NotFound(w, r)
			return
		}
		if errors.Is(err, repository.ErrCheckMetricsNotFound) {
			http.NotFound(w, r)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		"web/templates/base.html",
		"web/templates/check.html",
	)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	data := CheckPageData{Check: details.Check, Metrics: details.Metrics, Recommendations: details.Recommendations}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
