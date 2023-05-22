package models

type TemplateData struct {
	User         User
	Posts        []Post
	Post         Post
	Author       User
	Comments     []Comment
	PostReaction Reaction

	Error string
}
