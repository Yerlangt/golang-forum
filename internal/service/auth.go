package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"forum/internal/models"
	"forum/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// interface that will work with AuthService structure
type Auth interface {
	CreateUser(user models.User) error
	SetSession(username, password string) (models.Session, error)
	UserByToken(token string) (models.User, error)
	DeleteSession(token string) error
	GetUserByID(ID int) (models.User, error)
}

// structure that will return repository
type AuthService struct {
	repository repository.Auth
}

var (
	ErrNoUser        = errors.New("user doesn't exist")
	ErrWrongPassword = errors.New("incorrect password")
	ErrUserTaken     = errors.New("user with this data is taken")
)

func NewAuthService(repository repository.Auth) *AuthService {
	return &AuthService{
		repository: repository,
	}
}

// duration of the session will be one hour (from the log in moment)
const sessionTime = time.Hour

// creating user (with data validation), works with repository data
func (s *AuthService) CreateUser(user models.User) error {
	// check for uniqness of user's mail (using data from the bd)
	if _, err := s.repository.GetUser("", user.Email); err != sql.ErrNoRows {
		if err == nil {
			return ErrUserTaken
		}
		return err
	}

	// check for uniqness of user's username (using data from the bd)
	if _, err := s.repository.GetUser(user.UserName, ""); err != sql.ErrNoRows {
		if err == nil {
			return ErrUserTaken
		}
		return err
	}

	// validate email
	res, strErr := validateEmail(user.Email)
	if !res {
		return errors.New(strErr)
	}
	// validate password
	res, strErr = validatePassword(user.Password)
	if !res {
		return errors.New(strErr)
	}

	// validate username
	res, strErr = validateUsername(user.UserName)
	if !res {
		return errors.New(strErr)
	}

	// hashing password
	password, err := s.generateHashedPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = password

	return s.repository.CreateUser(user)
}

// setting session for user
func (s *AuthService) SetSession(username, password string) (models.Session, error) {
	user, err := s.checkPassword(username, password)
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
		return models.Session{}, err
	}

	return session, nil
}

// delete session by token
func (s *AuthService) DeleteSession(token string) error {
	return s.repository.DeleteSession(token)
}

// finding user by token in db
func (s *AuthService) UserByToken(token string) (models.User, error) {
	session, _ := s.repository.GetSessionByToken(token)
	user, _ := s.repository.GetUserById(session.UserID)

	return user, nil
}

// checking password with hashed on those in the db
func (s *AuthService) checkPassword(username, password string) (models.User, error) {
	user, err := s.repository.GetUser(username, "")
	if err != nil {
		return user, ErrNoUser
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, ErrWrongPassword
	}

	return user, nil
}

// generate random token for the user
func (s *AuthService) generateToken() (string, error) {
	const tokenLength = 32
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// using library bcrypt to hash password
func (s *AuthService) generateHashedPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash), err
}

func (s *AuthService) GetUserByID(ID int) (models.User, error) {
	return s.repository.GetUserById(ID)
}
