package repository

import (
	"database/sql"

	"forum/internal/models"
)

type Commentary interface {
	CreateComment(comment models.Comment) error
	GetCommentsByPostID(postID int) ([]models.Comment, error)
}

type CommentStorage struct {
	db *sql.DB
}

// return structure with db
func NewCommentStorage(db *sql.DB) *CommentStorage {
	return &CommentStorage{
		db: db,
	}
}

func (s *CommentStorage) CreateComment(comment models.Comment) error {
	query := `
        INSERT INTO COMMENTS(AuthorID, PostID, Content) VALUES ($1, $2, $3)
    `

	if _, err := s.db.Exec(query, comment.AuthorID, comment.PostID, comment.Content); err != nil {
		return err
	}
	return nil
}

func (s *CommentStorage) GetCommentsByPostID(postID int) ([]models.Comment, error) {
	query := `
		SELECT ID, AuthorID, PostID, Content FROM COMMENTS  WHERE PostID=?
	`

	rows, err := s.db.Query(query, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.AuthorID, &comment.PostID, &comment.Content); err != nil {
			return comments, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
