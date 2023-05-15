package repository

import "database/sql"

type Repository struct {
	Auth
	Post
	Commentary
	Reaction
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Auth:       NewAuthStorage(db),
		Post:       NewPostStorage(db),
		Commentary: NewCommentStorage(db),
		Reaction:   NewReactioStorage(db),
	}
}
