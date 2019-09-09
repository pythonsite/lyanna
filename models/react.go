package models

type ReactItem struct {
	BaseModel
	PostID int64
	ReactionType int64
}