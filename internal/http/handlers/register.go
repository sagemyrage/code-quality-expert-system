package handlers

import (
	"html/template"
	"net/http"
)

func RegisterPage(w http.ResponseWriter, r *http.Request) {
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
