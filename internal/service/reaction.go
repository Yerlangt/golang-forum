package service

import (
	"database/sql"
	"errors"

	"forum/internal/models"
	"forum/internal/repository"
)

type Reaction interface {
	CreateReaction(reaction models.Reaction) error
	GetReactionByPostID(PostID, AuthorID int) (models.Reaction, error)
	GetReactionCountByPostID(PostID int) (int, int, error)
	GetReactionCountByCommentID(CommentID int) (int, int, error)
	GetReactionByCommentID(CommentID, AuthorID int) (string, error)
}
type ReactionService struct {
	repository repository.Reaction
}

func NewReactionService(repository repository.Reaction) *ReactionService {
	return &ReactionService{
		repository: repository,
	}
}

func (s *ReactionService) GetReactionCountByPostID(PostID int) (int, int, error) {
	return s.repository.GetReactionCountByPostID(PostID)
}

func (s *ReactionService) GetReactionCountByCommentID(CommentID int) (int, int, error) {
	return s.repository.GetReactionCountByCommentID(CommentID)
}

func (s *ReactionService) CreateReaction(reaction models.Reaction) error {
	var exist models.Reaction
	var err error

	if reaction.CommentID != 0 && reaction.PostID != 0 {
		return errors.New("bad request")
	}
	if reaction.CommentID == 0 {
		exist, err = s.repository.GetReactionByPostID(reaction.PostID, reaction.AuthorID)
	} else if reaction.PostID == 0 {
		exist, err = s.repository.GetReactionByCommentID(reaction.CommentID, reaction.AuthorID)
	}

	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if exist == (models.Reaction{}) {
		return s.repository.CreateReaction(reaction)
	} else {
		if exist.Type == reaction.Type {
			exist.Type = "none"
		} else {
			exist.Type = reaction.Type
		}
		return s.repository.ChangeReaction(exist)
	}
}

func (s *ReactionService) GetReactionByPostID(PostID, AuthorID int) (models.Reaction, error) {
	return s.repository.GetReactionByPostID(PostID, AuthorID)
}

func (s *ReactionService) GetReactionByCommentID(CommentID, AuthorID int) (string, error) {
	reaction, err := s.repository.GetReactionByCommentID(CommentID, AuthorID)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if reaction == (models.Reaction{}) {
		return "none", nil
	}
	return reaction.Type, nil
}
