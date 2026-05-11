package handlers

import (
	"html/template"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(
		"web/templates/base.html",
		"web/templates/home.html",
	)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
