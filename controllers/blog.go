package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
	"lyanna/models"
	"lyanna/utils"
	"net/http"
	"strconv"
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
	pagination := utils.Pagination{
		CurrentPage:1,
		PerPage:models.Conf.General.PerPage,
		Total:len(posts),
	}
	var perPosts []*models.Post
	if models.Conf.General.PerPage > len(posts) {
		perPosts = posts
	} else {
		perPosts = posts[:models.Conf.General.PerPage]
	}

	c.HTML(http.StatusOK, "front/index.html",gin.H{
		"posts":perPosts,
		"pagination":&pagination,
	})
}

func Archives(c *gin.Context) {
	var ArchiveResult = make(map[string][]*models.Post )
	allArchives, _ := models.ListPostArchives()
	fmt.Println(allArchives)
	var year []string
	for _, v := range allArchives {
		posts := models.ListPostByArchive(v.Year)
		ArchiveResult[v.Year] = posts
		year = append(year,v.Year)
	}
	fmt.Println(year)
	c.HTML(http.StatusOK,"front/archives.html", gin.H{
		"ArchiveResult":ArchiveResult,
		"year":year,
	})
}

func ArchivesByYear(c *gin.Context) {
	year := c.Param("year")
	var ArchiveResult = make(map[string][]*models.Post )
	posts := models.ListPostByArchive(year)
	ArchiveResult[year] = posts
	c.HTML(http.StatusOK,"front/archives.html",gin.H{
		"ArchiveResult":ArchiveResult,
	})

}

func AboutMe(c *gin.Context) {
	slug := c.Param("aboutme")
	if slug != "aboutme" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Println(slug)
	post,err := models.GetPostBySlug(slug)
	if err != nil {

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Println(post)
	tags ,_ := models.ListTagByPostID(post.ID)
	post.Tags = tags
	content := models.GetContent(int(post.ID))
	gitHubUser, _ := c.Get(models.CONTEXT_USER_KEY)
	policy := bluemonday.UGCPolicy()
	unsafe := blackfriday.MarkdownCommon([]byte(content))
	contentHtml:=template.HTML(string(policy.SanitizeBytes(unsafe)))
	c.HTML(http.StatusOK,"front/post.html",gin.H{
		"Post":post,
		"contentHtml":contentHtml,
		"Githubuser":gitHubUser,
	})

}

func GetSearch(c *gin.Context) {
	c.HTML(http.StatusOK,"front/search.html",nil)
}

func PostSearch(c *gin.Context)  {
	var (
		posts []*models.Post
		err error
	)
	posts , err = models.ListPublishedPost("")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var ret []map[string]interface{}
	for _, post := range posts {
		var Posts = make(map[string]interface{},1)
		post.Tags,_ = models.ListTagByPostID(post.ID)
		Posts["url"] = post.Url()
		Posts["tags"] = post.GetTagsArray()
		Posts["title"] = post.Title
		Posts["content"] = models.GetContent(int(post.ID))
		ret = append(ret,Posts)

	}

	c.JSON(http.StatusOK,ret)

}

func PostPage(c *gin.Context) {
	var (
		posts []*models.Post
		err error
	)
	page := c.Param("page")
	pageInt, _ := strconv.ParseInt(page,10,32)
	posts , err = models.ListPublishedPost("")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	for _, post := range posts {
		post.Tags,_ = models.ListTagByPostID(post.ID)
	}
	pagination := utils.Pagination{
		CurrentPage:int(pageInt),
		PerPage:models.Conf.General.PerPage,
		Total:len(posts),
	}
	start := (int(pageInt) -1) * models.Conf.General.PerPage
	var end int
	if start+models.Conf.General.PerPage > len(posts) {
		end = len(posts)
	} else {
		end = start+models.Conf.General.PerPage
	}
	perPosts := posts[start:end]
	c.HTML(http.StatusOK, "front/index.html",gin.H{
		"posts":perPosts,
		"pagination":&pagination,
	})
}