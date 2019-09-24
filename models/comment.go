package models

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
)

var RedisCommentKey string = "comments/%d/props/content"

type Comment struct {
	BaseModel
	GitHubID int64
	PostID int64
	Content string `gorm:"type:longtext"`
	RefID int64
}

func (comment *Comment) Insert()error{
	return 	DB.Create(comment).Error
}


func ListCommentsByPostID(postid int)([]*Comment, error){
	var comments []*Comment
	err := DB.Model(&Comment{}).Order("id desc").Find(&comments,"post_id=?",postid).Error
	return comments,err
}

func (comment *Comment) GitUser() *GitHubUser{
	gitUser,_ := GetGitUserByGid(comment.GitHubID)
	return gitUser
}

func (comment *Comment) CommentHTML() template.HTML {
	policy := bluemonday.UGCPolicy()
	unsafe := blackfriday.MarkdownCommon([]byte(comment.Content))
	contentHtml:=template.HTML(string(policy.SanitizeBytes(unsafe)))
	return contentHtml
}

func (comment *Comment) GetComment(commentID interface{})string {
	res := GetCommentContent(commentID)
	return res
}

func CommentCreatAndGetID(comment *Comment)error {
	err := DB.Create(comment).Row().Scan(comment)
	return err
}

