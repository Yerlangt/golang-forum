package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type Reaction interface {
	CreateReaction(reaction models.Reaction) error
	GetReactionByIDs(PostID, AuthorID int) (models.Reaction, error)
}
type ReactionService struct {
	repository repository.Reaction
}

func NewReactionService(repository repository.Reaction) *ReactionService {
	return &ReactionService{
		repository: repository,
	}
}

func (s *ReactionService) CreateReaction(reaction models.Reaction) error {
	exist, err := s.repository.GetReactionByIDs(reaction.PostID, reaction.AuthorID)
	if err != nil {
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

func (s *ReactionService) GetReactionByIDs(PostID, AuthorID int) (models.Reaction, error) {
	return s.repository.GetReactionByIDs(PostID, AuthorID)
}
