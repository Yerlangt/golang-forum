package models

type Reaction struct {
	ID        int
	PostID    int
	CommentID int
	AuthorID  int
	Type      string
}
