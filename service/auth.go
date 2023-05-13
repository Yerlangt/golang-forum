package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"forum/internal/models"
	"forum/internal/repository"
)

// interface that will work with AuthService structure
type ServiceAuth interface {
	CreateUser(user models.User) error
	SetSession(username, password string) (models.Session, error)
	UserByToken(token string) (models.User, error)
	DeleteSession(token string) error
}

// structure that will return repository
type AuthService struct {
	repository repository.Auth
}

// duration of the session will be one hour (from the log in moment)
const sessionTime = time.Hour

func NewAuthService(repository repository.Auth) *AuthService {
	return &AuthService{
		repository: repository,
	}
}

// creating user (with data validation), works with repository data
func (s *AuthService) CreateUser(user models.User) error {
	if _, err := s.repository.GetUser("", user.Email); err != sql.ErrNoRows {
		if err == nil {
			return errors.New("email address is already taken")
		}
		return err
	}

	if _, err := s.repository.GetUser(user.UserName, ""); err != sql.ErrNoRows {
		if err == nil {
			return errors.New("username is already taken")
		}
		return err
	}
	// validate email, password, username
	// hash password

	return s.repository.CreateUser(user)
}

func (s *AuthService) SetSession(username, password string) (models.Session, error) {
	user, err := s.repository.GetUser(username, "")
	if err != nil {
		return models.Session{}, err
	}
	s.repository.DeleteSessionById(user.ID)
	token, _ := s.generateToken()

	session := models.Session{
		UserID:         user.ID,
		Token:          token,
		ExpirationDate: time.Now().Add(sessionTime),
	}

	if err := s.repository.CreateSession(session); err != nil {
		fmt.Println("HERE ERROR creation session service")
		return models.Session{}, err
	}

	return session, nil
}

func (s *AuthService) DeleteSession(token string) error {
	return s.repository.DeleteSession(token)
}

func (s *AuthService) UserByToken(token string) (models.User, error) {
	session, _ := s.repository.GetSessionByToken(token)
	user, _ := s.repository.GetUserById(session.UserID)

	return user, nil
}

func (s *AuthService) generateToken() (string, error) {
	const tokenLength = 32
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
