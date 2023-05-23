package handler

import (
	"database/sql"
	"html/template"
	"net/http"

	"forum/internal/models"
)

var userTemp, userParse = template.ParseFiles("web/template/user_page.html")

func (h *Handler) userPage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKey).(models.User)

	if r.Method == http.MethodGet {

		posts, err := h.services.UserPage.GetPostsByID(user.ID)
		if err != nil && err != sql.ErrNoRows {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
		for i := range posts {
			likes, dislikes, err := h.services.Reaction.GetReactionCountByPostID(posts[i].ID)
			if err != nil && err != sql.ErrNoRows {
				h.ErrorPage(w, http.StatusInternalServerError, err)
				return
			} else {
				posts[i].LikeCount = likes
				posts[i].DislikeCount = dislikes
			}
			commentCount, err := h.services.Commentary.GetCommentCountByPostID(posts[i].ID)
			if err != nil && err != sql.ErrNoRows {
				h.ErrorPage(w, http.StatusInternalServerError, err)
				return
			} else {
				posts[i].CommentCount = commentCount
			}
			categories, err := h.services.GetCategoriesByPostId(posts[i].ID)
			if err != nil && err != sql.ErrNoRows {
				h.ErrorPage(w, http.StatusInternalServerError, err)
				return
			} else {
				posts[i].Category = categories
			}
			posts[i].ShortVersion = h.services.GetShortVersionContent(posts[i].Content)

		}
		data := models.TemplateData{
			User:  user,
			Posts: posts,
		}
		if err := userTemp.Execute(w, data); err != nil || userParse != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		h.ErrorPage(w, http.StatusMethodNotAllowed, nil)
	}
}
