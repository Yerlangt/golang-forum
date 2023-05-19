package models

type Comment struct {
	ID           int
	AuthorID     int
	PostID       int
	Content      string
	AuthorName   string
	Reaction     string
	LikeCount    int
	DislikeCount int
}
