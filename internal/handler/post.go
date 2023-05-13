package handler

import (
	"net/http"
	"text/template"

	"forum/internal/models"
)

var postTemp, postParse = template.ParseFiles("templates/create_post.html")

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if err := postTemp.Execute(w, nil); err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
	} else if r.Method == http.MethodPost {
		user := r.Context().Value(ctxKey).(models.User)
		if err := r.ParseForm(); err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}

		title, err1 := r.Form["title"]
		content, err2 := r.Form["content"]

		if !err1 || !err2 {
			h.ErrorPage(w, http.StatusBadRequest, nil)
			return
		}

		post := models.Post{
			Title:    title[0],
			Content:  content[0],
			AuthorID: user.ID,
		}

		if err := h.services.ServicePost.CreatePost(post); err != nil {
			// error out of Validation
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		h.ErrorPage(w, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) postPage(w http.ResponseWriter, r *http.Request) {
}
