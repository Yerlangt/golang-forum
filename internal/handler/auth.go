package handler

import (
	"errors"
	"html/template"
	"net/http"
	"time"

	"forum/internal/models"
	"forum/internal/service"
)

var (
	signup, signupParse = template.ParseFiles("web/template/signup.html")
	signin, signinParse = template.ParseFiles("web/template/signin.html")
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if err := signup.Execute(w, nil); err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}

		email, err1 := r.Form["email"]
		username, err2 := r.Form["username"]
		password, err3 := r.Form["password"]

		if !err1 || !err2 || !err3 {
			h.ErrorPage(w, http.StatusBadRequest, nil)
			return
		}
		user := models.User{
			Email:    email[0],
			UserName: username[0],
			Password: password[0],
		}

		if err := h.services.Auth.CreateUser(user); err != nil {
			// error out of Validation
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		h.ErrorPage(w, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if err := signin.Execute(w, nil); err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}

		username, err1 := r.Form["username"]
		password, err2 := r.Form["password"]

		if !err1 || !err2 {
			h.ErrorPage(w, http.StatusBadRequest, nil)
			return
		}

		session, err := h.services.Auth.SetSession(username[0], password[0])
		if err != nil {
			if errors.Is(err, service.ErrNoUser) || errors.Is(err, service.ErrWrongPassword) {
				h.ErrorPage(w, http.StatusUnauthorized, err)
				return
			}
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
		cookie := &http.Cookie{
			Name:    "session_token",
			Value:   session.Token,
			Path:    "/",
			Expires: session.ExpirationDate,
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		h.ErrorPage(w, http.StatusMethodNotAllowed, nil)
		return
	}
}

func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			h.ErrorPage(w, http.StatusUnauthorized, err)
			return
		}
		h.ErrorPage(w, http.StatusBadRequest, err)
		return
	}
	if err := h.services.DeleteSession(cookie.Value); err != nil {
		h.ErrorPage(w, http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
