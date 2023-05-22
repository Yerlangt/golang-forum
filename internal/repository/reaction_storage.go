package repository

import (
	"database/sql"

	"forum/internal/models"
)

type Reaction interface {
	CreateReaction(reaction models.Reaction) error
	GetReactionByPostID(PostID, AuthorID int) (models.Reaction, error)
	GetReactionByCommentID(CommentID, AuthorID int) (models.Reaction, error)
	ChangeReaction(reaction models.Reaction) error
	GetReactionCountByPostID(PostID int) (int, int, error)
	GetReactionCountByCommentID(CommentID int) (int, int, error)
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

func (s *ReactioStorage) GetReactionByPostID(PostID, AuthorID int) (models.Reaction, error) {
	query := `
	SELECT ID, PostID, CommentID, AuthorID, Type FROM REACTIONS WHERE PostID=$1 AND AuthorID =$2;
	`

	var reaction models.Reaction

	if err := s.db.QueryRow(query, PostID, AuthorID).Scan(&reaction.ID, &reaction.PostID, &reaction.CommentID, &reaction.AuthorID, &reaction.Type); err != nil {
		return reaction, err
	}

	return reaction, nil
}

func (s *ReactioStorage) GetReactionByCommentID(CommentID, AuthorID int) (models.Reaction, error) {
	query := `
	SELECT ID, PostID, CommentID, AuthorID, Type FROM REACTIONS WHERE CommentID=$1 AND AuthorID =$2;
	`

	var reaction models.Reaction

	if err := s.db.QueryRow(query, CommentID, AuthorID).Scan(&reaction.ID, &reaction.PostID, &reaction.CommentID, &reaction.AuthorID, &reaction.Type); err != nil {
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

func (s *ReactioStorage) GetReactionCountByPostID(PostID int) (int, int, error) {
	queryDislikes := `
	SELECT COUNT(ID) FROM REACTIONS WHERE PostID=$1 AND Type ="dislike";
	`

	queryLikes := `
	SELECT COUNT(ID) FROM REACTIONS WHERE PostID=$1 AND Type ="like";
	`
	var likes, dislikes int

	if err := s.db.QueryRow(queryLikes, PostID).Scan(&likes); err != nil {
		return 0, 0, err
	}

	if err := s.db.QueryRow(queryDislikes, PostID).Scan(&dislikes); err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}

func (s *ReactioStorage) GetReactionCountByCommentID(CommentID int) (int, int, error) {
	queryDislikes := `
	SELECT COUNT(ID) FROM REACTIONS WHERE CommentID=$1 AND Type ="dislike";
	`

	queryLikes := `
	SELECT COUNT(ID) FROM REACTIONS WHERE CommentID=$1 AND Type ="like";
	`
	var likes, dislikes int

	if err := s.db.QueryRow(queryLikes, CommentID).Scan(&likes); err != nil {
		return 0, 0, err
	}

	if err := s.db.QueryRow(queryDislikes, CommentID).Scan(&dislikes); err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}
