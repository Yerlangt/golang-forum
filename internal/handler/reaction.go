package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/models"
)

func (h *Handler) createReaction(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKey).(models.User)
	postID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/posts/reaction/"))
	if err != nil {
		h.ErrorPage(w, http.StatusInternalServerError, err)
		return
	}
	if err := r.ParseForm(); err != nil {
		h.ErrorPage(w, http.StatusInternalServerError, err)
		return
	}

	var reactionType string
	if val, ok := r.Form["like"]; ok {
		reactionType = val[0]
	} else if val, ok := r.Form["dislike"]; ok {
		reactionType = val[0]
	} else {
		http.Redirect(w, r, fmt.Sprintf("/posts/%d", postID), http.StatusSeeOther)
		return
	}
	reaction := models.Reaction{
		PostID:   postID,
		AuthorID: user.ID,
		Type:     reactionType,
	}
	if err := h.services.Reaction.CreateReaction(reaction); err != nil {
		h.ErrorPage(w, http.StatusInternalServerError, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/posts/%d", postID), http.StatusSeeOther)
}
