package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"lyanna/models"
	"net/http"
	"strconv"
)

func PostIndex(c *gin.Context) {
	posts, _ := models.ListPosts()
	for _, post := range posts {
		tags, _:= models.ListTagByPostID(post.ID)
		post.Tags = tags
		log.Println(tags)
	}
	c.HTML(http.StatusOK, "admin/list_post.html", gin.H{
		"posts": posts,
		"post_count": len(posts),
	})
}

func GetEditPost(c *gin.Context) {
	id := c.Param("id")
	postID, _ := strconv.ParseUint(id,10,64)
	post, _ := models.GetPostByID(postID)
	tags ,_ := models.ListTagByPostID(post.ID)
	fmt.Println(tags)
	allTags, _ := models.ListALlTags()
	users, _:= models.ListUsers()
	post.Tags = tags
	var postTags []string
	for _,v:=range post.Tags {
		postTags = append(postTags, v.Name)
	}
	if post == nil {
		c.Redirect(http.StatusMovedPermanently,"/admin/users")
	}
	c.HTML(http.StatusOK, "admin/post.html",gin.H{
		"post":post,
		"users":users,
		"allTags":allTags,
		"postTags":postTags,
	})
}

func GetNewPost(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/post.html",nil)
}

func PostCreatePost(c *gin.Context) {
	id := c.Param("id")
	pID, _ := strconv.ParseInt(id,10,64)
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	summary := c.PostForm("summary")
	author := c.PostForm("author")
	authorID,_ := strconv.ParseInt(author,10,64)
	tags := c.PostFormArray("tags")
	fmt.Println("-----")
	fmt.Println(tags)
	canComment := "on" == c.PostForm("can_comment")
	publish := "on" == c.PostForm("publish")
	post := &models.Post{
		Title:title,
		Slug:slug,
		Summary:summary,
		AuthorID:int(authorID),
		CanComment:canComment,
		Published:publish,
	}
	post.ID = uint64(pID)
	post.Update()
	originPostTags,_ := models.ListTagByPostID(post.ID)
	originPostTagNames := models.GetTagNames(originPostTags)
	models.UpdateMultiTags(originPostTagNames, tags, int(post.ID))
	//newTags, _ :=  models.ListTagByPostID(post.ID)
	//post.Tags = newTags
	posts, _ := models.ListPosts()
	for _, post := range posts {
		tags, _:= models.ListTagByPostID(post.ID)
		post.Tags = tags
	}
	c.HTML(http.StatusOK,"admin/list_post.html",gin.H{
		"posts": posts,
		"post_count": len(posts),
		"msg":"User was successfully created.",
	})
}

