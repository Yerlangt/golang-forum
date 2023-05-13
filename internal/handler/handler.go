package handler

import (
	"net/http"

	"forum/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		services: service,
	}
}

// initialization of the routes in the handler
func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.homePage)
	mux.HandleFunc("/sign-up", h.signUp)
	mux.HandleFunc("/sign-in", h.signIn)
	mux.HandleFunc("/logout", h.isAuth(h.logOut))

	mux.HandleFunc("/posts/create", h.isAuth(h.createPost))
	mux.HandleFunc("/posts", h.postPage)

	mux.Handle("/web/", http.StripPrefix("/web", http.FileServer(http.Dir("web/"))))
	return h.middleware(mux)
}
