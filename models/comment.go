package models

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
)

type Comment struct {
	BaseModel
	GitHubID int64
	PostID int64
	RefID int64
}

func ListCommentsByPostID(postid int)([]*Comment, error){
	var comments []*Comment
	err := DB.Model(&Comment{}).Find(&comments,"post_id=?",postid).Error
	return comments,err
}

func (comment *Comment) GitUser() *GitHubUser{
	gitUser,_ := GetGitUserByGid(comment.GitHubID)
	return gitUser
}

func (comment *Comment) CommentHTML(content string) template.HTML {
	policy := bluemonday.UGCPolicy()
	unsafe := blackfriday.Run([]byte(content))
	contentHtml:=template.HTML(string(policy.SanitizeBytes(unsafe)))
	return contentHtml
}