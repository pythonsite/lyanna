package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
	"lyanna/models"
	"lyanna/utils"
	"net/http"
	"strconv"
)

func CreateComment(c *gin.Context) {
	content := c.Request.PostFormValue("content")
	postIDStr := c.Param("id")
	postID,_ := strconv.ParseInt(postIDStr,10,64)
	session := sessions.Default(c)
	gid := session.Get(models.SESSION_KEY)
	comment := models.Comment{
		GitHubID:gid.(int64),
		PostID:postID,
		RefID:0,
	}
	_ = models.CommentCreatAndGetID(&comment)
	models.SetCommentContent(int64(comment.ID), content)
	commentHTML,_ := utils.RenderSingleCommnet(&comment)
	c.JSON(http.StatusOK,gin.H{
		"r":0,
		"html":commentHTML,
	})
}

func CommentMarkdown(c *gin.Context) {
	commentContent := c.Request.PostFormValue("text")
	policy := bluemonday.UGCPolicy()
	unsafe := blackfriday.Run([]byte(commentContent))
	commentHtml:= template.HTML(string(policy.SanitizeBytes(unsafe)))
	c.JSON(http.StatusOK, gin.H{
		"r":0,
		"text": commentHtml,
	})
}