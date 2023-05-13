package service

import "forum/internal/repository"

type Service struct {
	Auth
	Post
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		Auth: NewAuthService(repository.Auth),
		Post: NewPostService(repository.Post),
	}
}
