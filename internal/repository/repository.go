package repository

import "database/sql"

type Repository struct {
	Auth
	Post
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		// As all methods in interface are linked with methods in struct, we transfer struct for this field
		// Sending nterface for privacy
		Auth: NewAuthStorage(db),
		Post: NewPostStorage(db),
	}
}
