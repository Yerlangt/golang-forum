package repository

import (
	"database/sql"

	"forum/internal/models"
)

type Post interface {
	CreatePost(post models.Post) error
	GetIDByCategory(elem string) (int, error)
	CreateLink(postID int, categoryID int) error
	GetAllPost() ([]models.Post, error)
	GetPostById(PostID int) (models.Post, error)
	GetLastID() (int, error)
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
	_, err := s.db.Exec(query, post.AuthorID, post.Title, post.Content)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStorage) GetLastID() (int, error) {
	var lastInsertID int64
	err := s.db.QueryRow("SELECT id FROM POSTS ORDER BY id DESC LIMIT 1").Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}
	return int(lastInsertID), nil
}

func (s *PostStorage) GetIDByCategory(elem string) (int, error) {
	query := `
        SELECT ID FROM CATEGORIES WHERE Category = ?
    `
	var categoryID int
	if err := s.db.QueryRow(query, elem).Scan(&categoryID); err != nil {
		return 0, err
	}
	return categoryID, nil
}

func (s *PostStorage) CreateLink(postID int, categoryID int) error {
	query := `
        INSERT INTO CATEGORYLINK (CategoryID, PostID) values ($1, $2);
    `
	if _, err := s.db.Exec(query, categoryID, postID); err != nil {
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
