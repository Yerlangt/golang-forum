package models

type TemplateData struct {
	User     User
	Posts    []Post
	Post     Post
	Comments []Comment
	Reaction Reaction
}
