package service

import (
	"errors"
	"strings"

	"forum/internal/models"
	"forum/internal/repository"
)

type Post interface {
	CreatePost(post models.Post) error
	GetAllPosts() ([]models.Post, error)
	GetPostById(postID int) (models.Post, error)
	GetPostsByCategory(category []string) ([]models.Post, error)
	GetLikedPostsByUserID(UserID int) ([]models.Post, error)
	GetCategoriesByPostId(postID int) ([]string, error)
}

type PostService struct {
	repository repository.Post
}

func NewPostService(repository repository.Post) *PostService {
	return &PostService{
		repository: repository,
	}
}

var (
	ErrEmptyPost = errors.New("empty post")
	ErrNoPost    = errors.New("post is not found")
)

func (s *PostService) CreatePost(post models.Post) error {
	if strings.TrimSpace(post.Content) == "" {
		return ErrEmptyPost
	}
	if err := s.repository.CreatePost(post); err != nil {
		return err
	}
	postID, err := s.repository.GetLastID()
	if err != nil {
		return err
	}
	for _, elem := range post.Category {
		categoryID, err := s.repository.GetIDByCategory(elem)
		if err != nil {
			return err
		}

		if err := s.repository.CreateLink(postID, categoryID); err != nil {
			return err
		}
	}
	return nil
}

func (s *PostService) GetPostById(postID int) (models.Post, error) {
	post, err := s.repository.GetPostById(postID)
	if err != nil {
		return post, err
	}
	return post, nil
}

func (s *PostService) GetCategoriesByPostId(postID int) ([]string, error) {
	categoriesID, err := s.repository.GetCategoriesByPostID(postID)
	var Categories []string
	if err != nil {
		// fmt.Println("service/post/70", err)
		return Categories, err
	}
	for _, id := range categoriesID {
		category, err := s.repository.GetCategoryByID(id)
		if err != nil {
			return Categories, err
		}
		Categories = append(Categories, category)
	}
	return Categories, nil
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	return s.repository.GetAllPost()
}

func (s *PostService) GetLikedPostsByUserID(UserID int) ([]models.Post, error) {
	return s.repository.GetLikedPostsByUserID(UserID)
}

// need to fix and delete similar posts
func (s *PostService) GetPostsByCategory(category []string) ([]models.Post, error) {
	var posts []models.Post
	for _, elem := range category {
		categoryID, err := s.repository.GetIDByCategory(elem)
		if err != nil {
			return nil, err
		}
		postTemp, err := s.repository.GetPostsByCategoryID(categoryID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, postTemp...)
	}
	posts = removeDuplicates(posts)
	return posts, nil
}

func removeDuplicates(posts []models.Post) []models.Post {
	uniqueMap := make(map[int]bool)
	uniquePost := make([]models.Post, 0)

	for _, post := range posts {
		if !uniqueMap[post.ID] {
			uniqueMap[post.ID] = true
			uniquePost = append(uniquePost, post)
		}
	}
	return uniquePost
}
