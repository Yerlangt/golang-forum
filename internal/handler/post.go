package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/models"
)

var (
	createPostTemp, createPostParse = template.ParseFiles("web/template/create_post.html")
	postTemp, postParse             = template.ParseFiles("web/template/post.html")
)

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if err := createPostTemp.Execute(w, nil); err != nil || createPostParse != nil {
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

		if err := h.services.Post.CreatePost(post); err != nil {
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
	user := r.Context().Value(ctxKey).(models.User)
	postID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/posts/"))
	if err != nil {
		h.ErrorPage(w, http.StatusNotFound, err)
		return
	}

	if r.Method == http.MethodGet {
		post, err := h.services.GetPostById(postID)
		if err != nil {
			h.ErrorPage(w, http.StatusNotFound, err)
			return
		}
		data := models.TemplateData{
			User: user,
			Post: post,
		}
		fmt.Println(data)
		if err := postTemp.Execute(w, data); err != nil || postParse != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
	}
}
