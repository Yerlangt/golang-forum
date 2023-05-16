package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

var errPage, errParse = template.ParseFiles("web/template/error.html")

func (h *Handler) ErrorPage(w http.ResponseWriter, status int, err error) {
	fmt.Println(err)
	w.WriteHeader(status)
	if errParse == nil {
		if err := errPage.Execute(w, status); err == nil {
			return
		}
	}
	http.Error(w, errParse.Error(), http.StatusInternalServerError)
}
