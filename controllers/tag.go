package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"lyanna/models"
	"net/http"
	"strconv"
)

func Tags(c *gin.Context) {
	tags, err := models.ListTag()
	if err != nil {
		msg := fmt.Sprintf("list tag err:%v",err)
		Logger.Fatal(msg)
	}
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
	tagID, err := strconv.ParseInt(tagStr,10,64)
	if err != nil {
		msg := fmt.Sprintf("parse int err:%v",err)
		Logger.Fatal(msg)
	}
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
		post.Tags,err = models.ListTagByPostID(post.ID)
		if err != nil {
			msg := fmt.Sprintf("list tag by postID err:%v",err)
			Logger.Error(msg)
			continue
		}
	}
	c.HTML(http.StatusOK, "front/tag.html",gin.H{
		"posts":posts,
		"tagName":tagName,
	})



}