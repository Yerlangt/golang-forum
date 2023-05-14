package service

import (
	"errors"
	"fmt"
	"strings"

	"forum/internal/models"
	"forum/internal/repository"
)

type Post interface {
	CreatePost(post models.Post) error
	GetAllPosts() ([]models.Post, error)
	GetPostById(postID int) (models.Post, error)
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
	postID, err := s.repository.GetLastID()
	fmt.Println("service/post/46: ", err)
	if err != nil {
		return err
	}
	for _, elem := range post.Category {
		categoryID, err := s.repository.GetIDByCategory(elem)
		if err != nil {
			return err
		}
		fmt.Println("service/post/48: ", postID, categoryID)
		if err := s.repository.CreateLink(postID, categoryID); err != nil {
			return err
		}
	}
	return nil
}

func (s *PostService) GetPostById(postID int) (models.Post, error) {
	post, err := s.repository.GetPostById(postID)
	if err != nil {
		return post, err
	}
	return post, nil
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	return s.repository.GetAllPost()
}
