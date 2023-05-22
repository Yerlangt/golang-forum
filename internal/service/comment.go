package service

import (
	"errors"
	"strings"

	"forum/internal/models"
	"forum/internal/repository"
)

type Commentary interface {
	CreateComment(comment models.Comment) error
	GetCommentsByPostID(postID int) ([]models.Comment, error)
	GetCommentCountByPostID(PostID int) (int, error)
}

type CommentService struct {
	repository repository.Commentary
}

func NewCommentService(repository repository.Commentary) *CommentService {
	return &CommentService{
		repository: repository,
	}
}

func (s *CommentService) CreateComment(comment models.Comment) error {
	if strings.TrimSpace(comment.Content) == "" {
		return errors.New("empty comment")
	}
	return s.repository.CreateComment(comment)
}

func (s *CommentService) GetCommentsByPostID(postID int) ([]models.Comment, error) {
	return s.repository.GetCommentsByPostID(postID)
}

func (s *CommentService) GetCommentCountByPostID(PostID int) (int, error) {
	return s.repository.GetCommentCountByPostID(PostID)
}
