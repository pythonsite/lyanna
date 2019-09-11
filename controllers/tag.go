package controllers

import (
	"github.com/gin-gonic/gin"
	"lyanna/models"
	"net/http"
	"strconv"
)

func Tags(c *gin.Context) {
	tags, _ := models.ListTag()
	c.HTML(http.StatusOK,"front/tags.html",gin.H{
		"tags":tags,
	})
}

func Tag(c *gin.Context) {
	var (
		posts []*models.Post
		//total int
		//policy    *bluemonday.Policy
		err error
	)
	tagStr := c.Param("id")
	tagID, _ := strconv.ParseInt(tagStr,10,64)
	tagName := models.GetTagNameByID(int(tagID))
	posts , err = models.ListPublishedPost(tagStr)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	//total, err = models.CountPostByTag(tagName)
	//if err != nil {
	//	c.AbortWithStatus(http.StatusInternalServerError)
	//	return
	//}
	//policy = bluemonday.StrictPolicy()
	for _, post := range posts {
		post.Tags,_ = models.ListTagByPostID(post.ID)
	}
	c.HTML(http.StatusOK, "front/tag.html",gin.H{
		"posts":posts,
		"tagName":tagName,
	})



}