package controllers

import (
	"github.com/gin-gonic/gin"
	"lyanna/models"
	"net/http"
)

func PostPublish(c *gin.Context) {
	postID := c.Param("postID")
	var H = gin.H{}
	post, _ := models.GetPostByID(postID)
	post.Published = true
	H["r"] = 0
	post.Update()
	c.JSON(http.StatusOK,H)
}

func DeletePublish(c *gin.Context) {
	postID := c.Param("postID")
	var H = gin.H{}
	post, _ := models.GetPostByID(postID)
	post.Published = false
	H["r"] = 0
	post.Update()
	c.JSON(http.StatusOK,H)
}