package handler

import (
	"net/http"
	"text/template"

	"forum/internal/models"
)

var index, indParse = template.ParseFiles("templates/index.html")

// home page with path "/"
func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.ErrorPage(w, http.StatusNotFound, nil)
		return
	}
	user := r.Context().Value(ctxKey).(models.User)

	posts, err := h.services.GetAllPosts()
	if err != nil {
		h.ErrorPage(w, http.StatusInternalServerError, err)
		return
	}
	data := models.TemplateData{
		User:  user,
		Posts: posts,
	}
	index.Execute(w, data)
}
