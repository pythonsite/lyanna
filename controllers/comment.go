package controllers

import (
	"fmt"
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

func Comments(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, _ := strconv.ParseInt(postIDStr,10,64)
	post,err := models.GetPostByID(postID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"r":1,
			"msg": "Post not exist",
		})
		return
	}
	pageStr := c.DefaultQuery("page","1")
	perPageStr := c.DefaultQuery("per_page","10")
	page, _ := strconv.ParseInt(pageStr,10,64)
	perPage,_ := strconv.ParseInt(perPageStr,10,64)
	comments, _ := models.ListCommentsByPostID(int(postID))
	start := (page - 1) * perPage
	var end int64
	if (start+ perPage) > int64(len(comments)) -1 {
		end = int64(len(comments))
	} else {
		end = start+ perPage
	}
	comments = comments[start: end]
	var pages int
	if len(comments) % 10 == 0 {
		pages = len(comments) / 10
	} else {
		pages =  len(comments) / 10 + 1
	}
	gitHubUser, _ := c.Get(models.CONTEXT_USER_KEY)
 	hh := utils.HH {
 		comments,
 		gitHubUser,
 		post,
		pages,
		len(comments),
	}
	fmt.Println(456)
	commentsHTML,_:= utils.RenderAllComment(hh)
	c.JSON(http.StatusOK,gin.H{
		"r":0,
		"html":commentsHTML,
	})

	//c.JSON(http.StatusOK,gin.H{
	//	"r": 0,
	//	"comments":comments,
	//	"githubuser":gitHubUser,
	//	"post":post,
	//})
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