package models

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
)

var RedisCommentKey string = "comments/%d/props/content"

type Comment struct {
	BaseModel
	GitHubID int64
	PostID int64
	RefID int64
}

func (comment *Comment) Insert()error{
	return 	DB.Create(comment).Error
}


func ListCommentsByPostID(postid int)([]*Comment, error){
	var comments []*Comment
	err := DB.Model(&Comment{}).Order("id desc").Find(&comments,"post_id=?",postid).Error
	fmt.Println(comments)
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

func (comment *Comment) GetComment(commentID interface{})string {
	res := GetCommentContent(commentID)
	return res
}

func CommentCreatAndGetID(comment *Comment)error {
	err := DB.Create(comment).Row().Scan(comment)
	return err
}

