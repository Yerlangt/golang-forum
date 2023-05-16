package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"reflect"

	"forum/internal/models"
)

var index, indParse = template.ParseFiles("web/template/index.html")

// home page with path "/"
func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || indParse != nil {
		h.ErrorPage(w, http.StatusNotFound, nil)
		return
	}
	user := r.Context().Value(ctxKey).(models.User)
	categories := memberTest(r, "ctgr")
	fmt.Println(categories)
	var posts []models.Post
	var err error
	if reflect.DeepEqual(categories, []string{"news", "sport", "music", "kids", "hobbies", "programming", "art", "cooking", "other"}) {
		posts, err = h.services.GetAllPosts()
	} else {
		posts, err = h.services.GetPostsByCategory(categories)
	}
	fmt.Println(err)
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
