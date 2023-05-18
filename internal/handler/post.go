package handler

import (
	"html/template"
	"log"
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
		category, err3 := r.Form["category1"]
		if category == nil {
			category = append(category, "other")
			err3 = true
		}

		if !err1 || !err2 || !err3 {
			h.ErrorPage(w, http.StatusBadRequest, nil)
			return
		}

		post := models.Post{
			Title:    title[0],
			Content:  content[0],
			AuthorID: user.ID,
			Category: category,
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
		author, err := h.services.Auth.GetUserByID(post.AuthorID)
		if err != nil {
			log.Printf("error getting author by ID: %s", err)
		}
		comments, err := h.services.GetCommentsByPostID(postID)
		if err != nil {
			log.Printf("error getting comments by post ID: %s", err)
		}
		reaction, err := h.services.Reaction.GetReactionByIDs(postID, user.ID)
		if err != nil {
			log.Printf("error getting GetReactionByIDs: %s", err)
		}
		likes, dislikes, err := h.services.Reaction.GetReactionCountByPostID(postID)
		if err != nil {
			log.Printf("error getting GetReactionCount: %s", err)
		} else {
			post.LikeCount = likes
			post.DislikeCount = dislikes
		}

		data := models.TemplateData{
			User:     user,
			Post:     post,
			Comments: comments,
			Reaction: reaction,
			Author:   author,
		}

		if err := postTemp.Execute(w, data); err != nil || postParse != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}

		content, ok := r.Form["content"]
		if !ok {
			h.ErrorPage(w, http.StatusBadRequest, nil)
			return
		}
		comment := models.Comment{
			AuthorID:   user.ID,
			PostID:     postID,
			Content:    content[0],
			AuthorName: user.UserName,
		}
		if err := h.services.Commentary.CreateComment(comment); err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)

	} else {
		h.ErrorPage(w, http.StatusMethodNotAllowed, nil)
	}
}
