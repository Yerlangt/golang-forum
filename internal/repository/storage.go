package repository

import (
	"database/sql"
)

func OpenDB(dbname string) (*sql.DB, error) {
	// open or create db with "dbname"
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, err
	}
	// Ping verifies a connection to the database is still alive, establishing a connection if necessary.
	if err = db.Ping(); err != nil {
		return nil, err
	}
	if err := initTables(db); err != nil {
		return nil, err
	}
	return db, nil
}

// initalization of the tables in the db
func initTables(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS USERS(
			ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			Username TEXT NOT NULL UNIQUE,
			Email TEXT NOT NULL UNIQUE,
			Password TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS SESSIONS(
			ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			UserID INTEGER NOT NULL UNIQUE,
			Token VARCHAR(32) NOT NULL ,
			ExpireDate DATATIME NOT NULL
		);
		CREATE TABLE IF NOT EXISTS POSTS(
			ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			AuthorID INTEGER NOT NULL,
			Title TEXT NOT NULL,
			Content TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS COMMENTS(
			ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			AuthorID INTEGER NOT NULL,
			PostID INTEGER NOT NULL,
			Content TEXT NOT NULL
		)
	`
	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
