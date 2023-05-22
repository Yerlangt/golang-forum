package handler

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"

	"forum/internal/models"
)

var userTemp, userParse = template.ParseFiles("web/template/user_page.html")

func (h *Handler) userPage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKey).(models.User)
	userNick := strings.TrimPrefix(r.URL.Path, "/user/")
	if r.Method == http.MethodGet {
		userID, err := h.services.UserPage.GetUserIDByNickName(userNick)
		if err != nil {

			h.ErrorPage(w, http.StatusNotFound, err)
			return
		}
		posts, err := h.services.UserPage.GetPostsByID(userID)
		if err != nil {
			log.Printf("error getting posts by user ID: %s", err)
		}
		for i := range posts {
			likes, dislikes, err := h.services.Reaction.GetReactionCountByPostID(posts[i].ID)
			if err != nil {
				log.Printf("error getting GetReactionCount: %s", err)
			} else {
				posts[i].LikeCount = likes
				posts[i].DislikeCount = dislikes
			}
			commentCount, err := h.services.Commentary.GetCommentCountByPostID(posts[i].ID)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("error getting GetCommentCountByPostID: %s", err)
			} else {
				posts[i].CommentCount = commentCount
			}
			categories, err := h.services.GetCategoriesByPostId(posts[i].ID)
			if err != nil {
				log.Printf("error getting GetCategories: %s", err)
			} else {
				posts[i].Category = categories
			}
			if len(posts[i].Content) > 200 {
				shortV := posts[i].Content[:200]
				words := strings.Split(shortV, " ")
				words = words[:len(words)-1]
				shortV = strings.Join(words, " ")
				posts[i].ShortVersion = shortV + " ..."
			} else {
				posts[i].ShortVersion = posts[i].Content
			}
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
