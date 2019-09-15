package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
	"lyanna/models"
	"net/http"
	"strconv"
)

func PostIndex(c *gin.Context) {
	posts, _ := models.ListPosts()
	for _, post := range posts {
		tags, _:= models.ListTagByPostID(post.ID)
		post.Tags = tags
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
	content := models.GetContent(int(postID))
	c.HTML(http.StatusOK, "admin/post.html",gin.H{
		"post":post,
		"users":users,
		"allTags":allTags,
		"postTags":postTags,
		"content": content,
	})
}

func GetNewPost(c *gin.Context) {
	allTags, _ := models.ListALlTags()
	users, _:= models.ListUsers()
	c.HTML(http.StatusOK, "admin/post.html",gin.H{
		"allTags":allTags,
		"users":users,
	})
}

func AddPost(c *gin.Context){
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	summary := c.PostForm("summary")
	author := c.PostForm("author")
	authorID,_ := strconv.ParseInt(author,10,64)
	tags := c.PostFormArray("tags")
	content := c.PostForm("content")
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
	_ = models.PostCreatAndGetID(post)
	models.SetContent(int(post.ID), content)
	models.UpdateMultiTags([]string{}, tags, int(post.ID))
	posts, _ := models.ListPosts()
	for _, post := range posts {
		tags, _:= models.ListTagByPostID(post.ID)
		post.Tags = tags
	}
	c.HTML(http.StatusOK,"admin/list_post.html",gin.H{
		"posts": posts,
		"post_count": len(posts),
		"msg":"Post was successfully created.",
	})
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	pID, _ := strconv.ParseInt(id,10,64)
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	summary := c.PostForm("summary")
	author := c.PostForm("author")
	authorID,_ := strconv.ParseInt(author,10,64)
	tags := c.PostFormArray("tags")
	content := c.PostForm("content")
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
	models.SetContent(int(pID), content)
	originPostTags,_ := models.ListTagByPostID(post.ID)
	originPostTagNames := models.GetTagNames(originPostTags)
	models.UpdateMultiTags(originPostTagNames, tags, int(post.ID))
	posts, _ := models.ListPosts()
	for _, post := range posts {
		tags, _:= models.ListTagByPostID(post.ID)
		post.Tags = tags
	}
	c.HTML(http.StatusOK,"admin/list_post.html",gin.H{
		"posts": posts,
		"post_count": len(posts),
		"msg":"Update post successfully.",
	})
}

func GetPost(c *gin.Context) {
	id := c.Param("id")
	postID, _ := strconv.ParseUint(id,10,64)
	post, _ := models.GetPostByID(postID)
	tags ,_ := models.ListTagByPostID(post.ID)
	post.Tags = tags
	content := models.GetContent(int(postID))
	comments, _ := models.ListCommentsByPostID(int(postID))
	gitHubUser, _ := c.Get(models.CONTEXT_USER_KEY)

	policy := bluemonday.UGCPolicy()
	unsafe := blackfriday.Run([]byte(content))
	contentHtml:=template.HTML(string(policy.SanitizeBytes(unsafe)))
	c.HTML(http.StatusOK,"front/post.html",gin.H{
		"post":post,
		"contentHtml":contentHtml,
		"comments": comments,
		"githubuser":gitHubUser,
	})
}

