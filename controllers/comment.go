package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"lyanna/models"
	"net/http"
	"strconv"
)

func CreateComment(c *gin.Context) {
	content := c.Request.PostFormValue("content")
	postIDStr := c.Param("id")
	postID,_ := strconv.ParseInt(postIDStr,10,64)
	session := sessions.Default(c)
	gid := session.Get(models.SESSION_KEY)
	comment := &models.Comment{
		GitHubID:gid.(int64),
		PostID:postID,
		RefID:0,
	}
	_ = models.CommentCreatAndGetID(comment)
	models.SetCommentContent(int64(comment.ID), content)
	c.JSON(http.StatusOK,gin.H{
		"r":0,
	})

}
