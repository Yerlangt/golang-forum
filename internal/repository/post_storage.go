package repository

import (
	"database/sql"

	"forum/internal/models"
)

type Post interface {
	CreatePost(post models.Post) error
	GetAllPost() ([]models.Post, error)
	GetPostById(PostID int) (models.Post, error)
}

type PostStorage struct {
	db *sql.DB
}

func NewPostStorage(db *sql.DB) *PostStorage {
	return &PostStorage{
		db: db,
	}
}

func (s *PostStorage) CreatePost(post models.Post) error {
	query := `
        INSERT INTO POSTS (AuthorID, Title, Content) VALUES ($1, $2, $3)
    `

	if _, err := s.db.Exec(query, post.AuthorID, post.Title, post.Content); err != nil {
		return err
	}
	return nil
}

func (s *PostStorage) GetPostById(PostID int) (models.Post, error) {
	query := `
	SELECT ID, AuthorID, Title, Content FROM POSTS  WHERE ID=?;
	`

	var post models.Post

	if err := s.db.QueryRow(query, PostID).Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content); err != nil {
		return post, err
	}

	return post, nil
}

func (s *PostStorage) GetAllPost() ([]models.Post, error) {
	query := `
		SELECT ID, AuthorID, Title, Content FROM POSTS 
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content); err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
