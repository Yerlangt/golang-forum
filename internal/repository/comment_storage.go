package repository

import (
	"database/sql"

	"forum/internal/models"
)

type Commentary interface {
	CreateComment(comment models.Comment) error
	GetCommentsByPostID(postID int) ([]models.Comment, error)
	GetCommentCountByPostID(PostID int) (int, error)
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
        INSERT INTO COMMENTS(AuthorID, PostID, Content, AuthorName) VALUES ($1, $2, $3, $4)
    `

	if _, err := s.db.Exec(query, comment.AuthorID, comment.PostID, comment.Content, comment.AuthorName); err != nil {
		return err
	}

	return nil
}

func (s *CommentStorage) GetCommentsByPostID(postID int) ([]models.Comment, error) {
	query := `
		SELECT ID, AuthorID, PostID, Content, AuthorName FROM COMMENTS  WHERE PostID=?
	`

	rows, err := s.db.Query(query, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.AuthorID, &comment.PostID, &comment.Content, &comment.AuthorName); err != nil {
			return comments, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (s *CommentStorage) GetCommentCountByPostID(PostID int) (int, error) {
	query := `
	SELECT COUNT(ID) FROM COMMENTS WHERE PostID=?;
	`

	var comments int

	if err := s.db.QueryRow(query, PostID).Scan(&comments); err != nil {
		return 0, err
	}

	return comments, nil
}
