package handler

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"reflect"

	"forum/internal/models"
)

var index, indParse = template.ParseFiles("web/template/index.html")

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.ErrorPage(w, http.StatusNotFound, nil)
		return
	}
	if r.Method != http.MethodGet {
		h.ErrorPage(w, http.StatusMethodNotAllowed, errors.New("homepage: method not allowed"))
		return
	}

	categories := memberTest(r, "ctgr")

	var posts []models.Post
	var err error
	if reflect.DeepEqual(categories, []string{"news", "sport", "music", "kids", "hobbies", "programming", "art", "cooking", "other"}) {
		posts, err = h.services.GetAllPosts()
	} else {
		posts, err = h.services.GetPostsByCategory(categories)
	}

	if err != nil && err != sql.ErrNoRows {
		h.ErrorPage(w, http.StatusInternalServerError, err)
		return
	}

	for i := range posts {
		likes, dislikes, err := h.services.Reaction.GetReactionCountByPostID(posts[i].ID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error is on homepage: getting GetReactionCount: %s", err)
		} else {
			posts[i].LikeCount = likes
			posts[i].DislikeCount = dislikes
		}
		commentCount, err := h.services.Commentary.GetCommentCountByPostID(posts[i].ID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error is on homepage: getting GetCommentCount: %s", err)
		} else {
			posts[i].CommentCount = commentCount
		}
		categories, err := h.services.GetCategoriesByPostId(posts[i].ID)
		if err != nil {
			log.Printf("Error is on homepage: getting GetCategories %s", err)
		} else {
			posts[i].Category = categories
		}

		posts[i].ShortVersion = h.services.GetShortVersionContent(posts[i].Content)

	}
	user := r.Context().Value(ctxKey).(models.User)

	data := models.TemplateData{
		User:  user,
		Posts: posts,
	}
	if err := index.Execute(w, data); err != nil || indParse != nil {
		h.ErrorPage(w, http.StatusInternalServerError, err)
		return
	}
}
