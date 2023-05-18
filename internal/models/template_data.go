package models

type TemplateData struct {
	User     User
	Posts    []Post
	Post     Post
	Author   User
	Comments []Comment
	Reaction Reaction
	Error    string
}
