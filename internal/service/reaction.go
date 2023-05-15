package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type Reaction interface {
	CreateReaction(reaction models.Reaction) error
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
	return s.repository.CreateReaction(reaction)
}
