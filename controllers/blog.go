package controllers

import (
	"github.com/gin-gonic/gin"
	"lyanna/models"
	"net/http"
)

func Index(c *gin.Context) {
	var (
		posts []*models.Post
		//total int
		//policy    *bluemonday.Policy
		err error
	)
	posts , err = models.ListPublishedPost("")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	for _, post := range posts {
		post.Tags,_ = models.ListTagByPostID(post.ID)
	}
	c.HTML(http.StatusOK, "front/index.html",gin.H{
		"posts":posts,
	})
}
