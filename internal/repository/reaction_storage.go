package repository

import (
	"database/sql"

	"forum/internal/models"
)

type Reaction interface {
	CreateReaction(reaction models.Reaction) error
}

type ReactioStorage struct {
	db *sql.DB
}

// return structure with db
func NewReactioStorage(db *sql.DB) *ReactioStorage {
	return &ReactioStorage{
		db: db,
	}
}

func (s *ReactioStorage) CreateReaction(reaction models.Reaction) error {
	query := `
        INSERT INTO REACTIONS (PostID, CommentID, AuthorID, Type) VALUES ($1, $2, $3, $4)
    `
	_, err := s.db.Exec(query, reaction.PostID, reaction.CommentID, reaction.AuthorID, reaction.Type)
	if err != nil {
		return err
	}

	return nil
}
