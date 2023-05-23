package handler

import (
	"database/sql"
	"errors"
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
	user := r.Context().Value(ctxKey).(models.User)
	if r.Method == http.MethodGet {
		data := models.TemplateData{
			User: user,
		}
		if err := createPostTemp.Execute(w, data); err != nil || createPostParse != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
		title, err1 := r.Form["title"]
		content, err2 := r.Form["content"]
		category, err3 := r.Form["category1"]

		// ADD HERE CHECK FOR THE CATEGORIES NAMES!!!!
		if category == nil {
			category = append(category, "other")
			err3 = true
		}
		check := h.services.Post.CheckCategory(category)
		// fmt.Println("CHECK: ", check)
		if !err1 || !err2 || !err3 || !check {
			h.ErrorPage(w, http.StatusBadRequest, errors.New("error: status bad request"))
			return
		}
		post := models.Post{
			Title:    title[0],
			Content:  content[0],
			AuthorID: user.ID,
			Category: category,
		}

		if err := h.services.Post.CreatePost(post); err != nil {
			data := models.TemplateData{
				Error: err.Error(),
				User:  user,
			}
			w.WriteHeader(400)
			if err := createPostTemp.Execute(w, data); err != nil {
				h.ErrorPage(w, http.StatusInternalServerError, err)
			}
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
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
		categories, err := h.services.GetCategoriesByPostId(postID)
		if err != nil && err != sql.ErrNoRows {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		} else {
			post.Category = categories
		}
		comments, err := h.services.GetCommentsByPostID(postID)
		if err != nil && err != sql.ErrNoRows {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}

		commentCount, err := h.services.Commentary.GetCommentCountByPostID(postID)
		if err != nil && err != sql.ErrNoRows {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		} else {
			post.CommentCount = commentCount
		}

		for i := range comments {
			commentType, err := h.services.Reaction.GetReactionByCommentID(comments[i].ID, user.ID)
			if err != nil && err != sql.ErrNoRows {
				h.ErrorPage(w, http.StatusInternalServerError, err)
				return
			}
			likes, dislikes, err := h.services.Reaction.GetReactionCountByCommentID(comments[i].ID)
			if err != nil && err != sql.ErrNoRows {
				h.ErrorPage(w, http.StatusInternalServerError, err)
				return
			}
			comments[i].Reaction = commentType
			comments[i].LikeCount = likes
			comments[i].DislikeCount = dislikes
		}
		postReaction, err := h.services.Reaction.GetReactionByPostID(postID, user.ID)
		if err != nil && err != sql.ErrNoRows {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
		likes, dislikes, err := h.services.Reaction.GetReactionCountByPostID(postID)
		if err != nil && err != sql.ErrNoRows {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		} else {
			post.LikeCount = likes
			post.DislikeCount = dislikes
		}

		data := models.TemplateData{
			User:         user,
			Post:         post,
			Comments:     comments,
			PostReaction: postReaction,
			Author:       author,
		}

		if err := postTemp.Execute(w, data); err != nil || postParse != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
	} else if r.Method == http.MethodPost {
		_, err := h.services.GetPostById(postID)
		if err != nil {
			h.ErrorPage(w, http.StatusBadRequest, err)
			return
		}
		if user == (models.User{}) {
			http.Redirect(w, r, "/sign-up", http.StatusSeeOther)
			return
		}
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

var likedTemp, likedPostParse = template.ParseFiles("web/template/liked_post.html")

func (h *Handler) likedPostPage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKey).(models.User)

	if r.Method == http.MethodGet {
		posts, err := h.services.Post.GetLikedPostsByUserID(user.ID)
		if err != nil && err != sql.ErrNoRows {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
		data := models.TemplateData{
			User:  user,
			Posts: posts,
		}
		if err := likedTemp.Execute(w, data); err != nil || likedPostParse != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		h.ErrorPage(w, http.StatusMethodNotAllowed, errors.New("error: method not allowed"))
	}
}
