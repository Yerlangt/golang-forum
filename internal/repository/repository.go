package repository

import "database/sql"

type Repository struct {
	Auth
	Post
	Commentary
	Reaction
	User
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Auth:       NewAuthStorage(db),
		Post:       NewPostStorage(db),
		Commentary: NewCommentStorage(db),
		Reaction:   NewReactionStorage(db),
		User:       NewUserStorage(db),
	}
}
