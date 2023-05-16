package repository

import (
	"database/sql"
	"fmt"

	"forum/internal/models"
)

type User interface {
	GetUserIDByNickName(nickName string) (int, error)
	GetPostsIDByUserID(userID int) ([]models.Post, error)
}

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) GetUserIDByNickName(nickName string) (int, error) {
	query := `SELECT ID FROM USERS WHERE Username = ?`
	var ID int
	if err := s.db.QueryRow(query, nickName).Scan(&ID); err != nil {
		return 0, err
	}
	fmt.Println("repo/userdata: ", ID)
	return ID, nil
}

func (s *UserStorage) GetPostsIDByUserID(userID int) ([]models.Post, error) {
	query := `SELECT ID, AuthorID, Title, Content FROM POSTS WHERE AuthorID = ?`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
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
