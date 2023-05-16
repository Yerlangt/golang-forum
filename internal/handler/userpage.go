package handler

import (
	"fmt"
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
			fmt.Println(err)
			h.ErrorPage(w, http.StatusNotFound, err)
			return
		}
		posts, err := h.services.UserPage.GetPostsByID(userID)
		if err != nil {
			log.Printf("error getting posts by user ID: %s", err)
		}
		data := models.TemplateData{
			User:  user,
			Posts: posts,
		}
		if err := userTemp.Execute(w, data); err != nil || userParse != nil {
			fmt.Println("error in Parcing")
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		h.ErrorPage(w, http.StatusMethodNotAllowed, nil)
	}
}
