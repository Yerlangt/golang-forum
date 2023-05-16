package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type UserPage interface {
	GetUserIDByNickName(nickname string) (int, error)
	GetPostsByID(id int) ([]models.Post, error)
}

type UserService struct {
	repository repository.User
}

func NewUserService(repository repository.User) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) GetUserIDByNickName(nickname string) (int, error) {
	return s.repository.GetUserIDByNickName(nickname)
}

func (s *UserService) GetPostsByID(id int) ([]models.Post, error) {
	return s.repository.GetPostsIDByUserID(id)
}
