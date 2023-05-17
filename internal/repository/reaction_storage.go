package repository

import (
	"database/sql"

	"forum/internal/models"
)

type Reaction interface {
	CreateReaction(reaction models.Reaction) error
	GetReactionByIDs(PostID, AuthorID int) (models.Reaction, error)
	ChangeReaction(reaction models.Reaction) error
}

type ReactioStorage struct {
	db *sql.DB
}

// return structure with db
func NewReactionStorage(db *sql.DB) *ReactioStorage {
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

func (s *ReactioStorage) GetReactionByIDs(PostID, AuthorID int) (models.Reaction, error) {
	query := `
	SELECT ID, PostID, CommentID, AuthorID, Type FROM REACTIONS WHERE PostID=$1 AND AuthorID =$2;
	`

	var reaction models.Reaction

	if err := s.db.QueryRow(query, PostID, AuthorID).Scan(&reaction.ID, &reaction.PostID, &reaction.CommentID, &reaction.AuthorID, &reaction.Type); err != nil {
		return reaction, err
	}

	return reaction, nil
}

func (s *ReactioStorage) ChangeReaction(reaction models.Reaction) error {
	query := `
	UPDATE REACTIONS SET Type = $1 WHERE ID=$2; 
    `
	_, err := s.db.Exec(query, reaction.Type, reaction.ID)
	if err != nil {
		return err
	}

	return nil
}
