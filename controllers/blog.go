package controllers

import (
	"github.com/gin-gonic/gin"
	"lyanna/models"
	"net/http"
)

func Index(c *gin.Context) {
	var (
		posts []*models.Post
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

func Archives(c *gin.Context) {
	var ArchiveResult = make(map[string][]*models.Post )
	allArchives, _ := models.ListPostArchives()
	for _, v := range allArchives {
		posts := models.ListPostByArchive(v.Year)
		ArchiveResult[v.Year] = posts
	}
	c.HTML(http.StatusOK,"front/archives.html", gin.H{
		"ArchiveResult":ArchiveResult,
	})
}