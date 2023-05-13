package service

import "forum/internal/repository"

type Service struct {
	ServiceAuth
	ServicePost
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		ServiceAuth: NewAuthService(repository.Auth),
		ServicePost: NewPostService(repository.Post),
	}
}
