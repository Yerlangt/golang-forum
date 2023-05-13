package service

import (
	"errors"
	"strings"

	"forum/internal/models"
	"forum/internal/repository"
)

type Post interface {
	CreatePost(post models.Post) error
	GetAllPosts() ([]models.Post, error)
}

type PostService struct {
	repository repository.Post
}

func NewPostService(repository repository.Post) *PostService {
	return &PostService{
		repository: repository,
	}
}

var (
	ErrEmptyPost = errors.New("empty post")
	ErrNoPost    = errors.New("post is not found")
)

func (s *PostService) CreatePost(post models.Post) error {
	if strings.TrimSpace(post.Content) == "" {
		return ErrEmptyPost
	}
	if err := s.repository.CreatePost(post); err != nil {
		return err
	}
	return nil
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	return s.repository.GetAllPost()
}
