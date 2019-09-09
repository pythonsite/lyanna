package models

type Comment struct {
	BaseModel
	GitHubID int64
	PostID int64
	ReactionType int64
	RefID int64
}