package handler

import (
	"net/http"
	"text/template"
)

var errPage, errParse = template.ParseFiles("templates/error.html")

func (h *Handler) ErrorPage(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	if errParse == nil {
		if err := errPage.Execute(w, status); err == nil {
			return
		}
	}
	http.Error(w, errParse.Error(), http.StatusInternalServerError)
}
