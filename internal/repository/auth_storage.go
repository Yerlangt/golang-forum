package repository

import (
	"database/sql"

	// models saves structs
	"forum/internal/models"
)

type Auth interface {
	CreateUser(user models.User) error
	GetUserByUsernameOrEmail(username, email string) (models.User, error)
	CreateSession(session models.Session) error
	GetSessionByToken(token string) (models.Session, error)
	GetUserById(id int) (models.User, error)
	DeleteSessionById(userID int) error
	DeleteSessionByToken(token string) error
}

type AuthStorage struct {
	db *sql.DB
}

// return structure with db
func NewAuthStorage(db *sql.DB) *AuthStorage {
	return &AuthStorage{
		db: db,
	}
}

// Adding information regarding new user to the db Users
func (s *AuthStorage) CreateUser(user models.User) error {
	query := `
		INSERT INTO USERS (Username, Email, Password) VALUES ($1, $2, $3);
	`
	if _, err := s.db.Exec(query, user.UserName, user.Email, user.Password); err != nil {
		return err
	}

	return nil
}

// Returning user information (struct form) from the db Users according to username and email
func (s *AuthStorage) GetUserByUsernameOrEmail(username, email string) (models.User, error) {
	query := `
		SELECT ID, Username, Email, Password FROM USERS WHERE Username=$1 or Email = $2;
	`

	var user models.User

	if err := s.db.QueryRow(query, username, email).Scan(&user.ID, &user.UserName, &user.Email, &user.Password); err != nil {
		return user, err
	}

	return user, nil
}

// Create session for user (adding token and expire date to db table Sessions)
func (s *AuthStorage) CreateSession(session models.Session) error {
	query := `
		INSERT INTO SESSIONS (UserID, Token, ExpireDate) VALUES ($1, $2, $3);
	`

	if _, err := s.db.Exec(query, session.UserID, session.Token, session.ExpirationDate); err != nil {
		return err
	}

	return nil
}

// Delete session from db table Sessions according to user token
func (s *AuthStorage) DeleteSessionByToken(token string) error {
	query := `
		DELETE FROM SESSIONS WHERE Token=?;
	`

	if _, err := s.db.Exec(query, token); err != nil {
		return err
	}
	return nil
}

// Delete session from db table Sessions according to user id
func (s *AuthStorage) DeleteSessionById(userID int) error {
	query := `
		DELETE FROM SESSIONS WHERE UserID=?;
	`

	if _, err := s.db.Exec(query, userID); err != nil {
		return err
	}
	return nil
}

// returning information regarding user's session from db table Sessions according to token
func (s *AuthStorage) GetSessionByToken(token string) (models.Session, error) {
	query := `
		SELECT ID, UserID, Token, ExpireDate FROM SESSIONS WHERE Token=$1;
	`

	var session models.Session

	if err := s.db.QueryRow(query, token).Scan(&session.ID, &session.UserID, &session.Token, &session.ExpirationDate); err != nil {
		return session, err
	}

	return session, nil
}

// returning information regarding user's session from db table Sessions according to id
func (s *AuthStorage) GetUserById(id int) (models.User, error) {
	query := `
		SELECT ID, Username, Email, Password FROM USERS WHERE ID=$1;
	`

	var user models.User

	if err := s.db.QueryRow(query, id).Scan(&user.ID, &user.UserName, &user.Email, &user.Password); err != nil {
		return user, err
	}

	return user, nil
}
