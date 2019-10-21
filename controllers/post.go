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
var Logger = models.Logger
// blackfriday 配置
const (
	commonHtmlFlags = 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES |
		blackfriday.HTML_NOFOLLOW_LINKS

	commonExtensions = 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS
)

func PostIndex(c *gin.Context) {
	posts, err := models.ListPosts()
	if err != nil {
		msg := fmt.Sprintf("get all posts err:%v",err)
		Logger.Fatal(msg)
	}
	for _, post := range posts {
		tags, err:= models.ListTagByPostID(post.ID)
		if err != nil {
			msg := fmt.Sprintf("list tag by postID err:%v",err)
			Logger.Error(msg)
			continue
		}
		post.Tags = tags
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
	c.HTML(http.StatusOK, "admin/list_post.html", gin.H{
		"posts": perPosts,
		"post_count": len(posts),
		"pagination":&pagination,
	})
}

func AdminPostPage(c *gin.Context) {
	page := c.Param("page")
	pageInt, _ := strconv.ParseInt(page,10,32)
	posts, _ := models.ListPosts()
	for _, post := range posts {
		tags, err:= models.ListTagByPostID(post.ID)
		if err != nil {
			msg := fmt.Sprintf("list tag by postID err:%v",err)
			Logger.Error(msg)
			continue
		}
		post.Tags = tags
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
	c.HTML(http.StatusOK, "admin/list_post.html", gin.H{
		"posts": perPosts,
		"post_count": len(posts),
		"pagination":&pagination,
	})
}

func GetEditPost(c *gin.Context) {
	id := c.Param("id")
	postID, err := strconv.ParseUint(id,10,64)
	if err != nil {
		msg := fmt.Sprintf("parse uint error:%v",err)
		Logger.Fatal(msg)
	}
	post, err := models.GetPostByID(postID)
	if err != nil {
		msg := fmt.Sprintf("get post by ID error:%v",err)
		Logger.Fatal(msg)
	}
	tags ,err := models.ListTagByPostID(post.ID)
	if err != nil {
		msg := fmt.Sprintf("list tag by postID error:%v",err)
		Logger.Fatal(msg)
	}
	allTags, err := models.ListALlTags()
	if err != nil {
		msg := fmt.Sprintf("list all tags error:%v",err)
		Logger.Fatal(msg)
	}
	users, err:= models.ListUsers()
	if err != nil {
		msg := fmt.Sprintf("list users error:%v",err)
		Logger.Fatal(msg)
	}
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
	allTags, err := models.ListALlTags()
	if err != nil {
		msg := fmt.Sprintf("list all tags error:%v",err)
		Logger.Fatal(msg)
	}
	users, err:= models.ListUsers()
	if err != nil {
		msg := fmt.Sprintf("list users error:%v",err)
		Logger.Fatal(msg)
	}
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
	authorID,err := strconv.ParseInt(author,10,64)
	if err != nil {
		msg := fmt.Sprintf("ParseInt error:%v",err)
		Logger.Fatal(msg)
	}
	tags := c.PostFormArray("tags")
	content := c.PostForm("content")
	canComment := "on" == c.PostForm("can_comment")
	publish := "on" == c.PostForm("publish")
	post := &models.Post{
		Title:title,
		Slug:slug,
		Summary:summary,
		AuthorID:int(authorID),
		Content:content,
		CanComment:canComment,
		Published:publish,
	}
	err = models.PostCreatAndGetID(post)
	if err != nil {
		msg := fmt.Sprintf("PostCreatAndGetID error:%v",err)
		Logger.Fatal(msg)
	}
	models.UpdateMultiTags([]string{}, tags, int(post.ID))
	posts, err := models.ListPosts()
	if err != nil {
		msg := fmt.Sprintf("list posts error:%v",err)
		Logger.Fatal(msg)
	}
	for _, post := range posts {
		tags, err:= models.ListTagByPostID(post.ID)
		if err != nil {
			msg := fmt.Sprintf("list tag by postID error:%v",err)
			Logger.Error(msg)
			continue
		}
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
	pID, err := strconv.ParseInt(id,10,64)
	if err != nil {
		msg := fmt.Sprintf("parse int error:%v",err)
		Logger.Fatal(msg)
	}
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	summary := c.PostForm("summary")
	author := c.PostForm("author")
	authorID,err := strconv.ParseInt(author,10,64)
	if err != nil {
		msg := fmt.Sprintf("parse int error:%v",err)
		Logger.Fatal(msg)
	}
	tags := c.PostFormArray("tags")
	content := c.PostForm("content")
	canComment := "on" == c.PostForm("can_comment")
	publish := "on" == c.PostForm("publish")
	post := &models.Post{
		Title:title,
		Slug:slug,
		Summary:summary,
		AuthorID:int(authorID),
		Content:content,
		CanComment:canComment,
		Published:publish,
	}
	post.ID = uint64(pID)
	post.Update()
	originPostTags,err := models.ListTagByPostID(post.ID)
	if err != nil {
		msg := fmt.Sprintf("list tag by postID error:%v",err)
		Logger.Fatal(msg)
	}
	originPostTagNames := models.GetTagNames(originPostTags)
	models.UpdateMultiTags(originPostTagNames, tags, int(post.ID))
	posts, err := models.ListPosts()
	if err != nil {
		msg := fmt.Sprintf("list posts error:%v",err)
		Logger.Fatal(msg)
	}
	for _, post := range posts {
		tags, err:= models.ListTagByPostID(post.ID)
		if err != nil {
			msg := fmt.Sprintf("list tag by postID error:%v",err)
			Logger.Error(msg)
			continue
		}
		post.Tags = tags
	}
	c.HTML(http.StatusOK,"admin/list_post.html",gin.H{
		"posts": posts,
		"post_count": len(posts),
		"msg":"Update post successfully.",
	})
}

func PreviewGetPost(c *gin.Context) {
	getPost(c, false)
}

func GetPost(c *gin.Context){
	getPost(c, true)
}

func getPost(c *gin.Context, isPublish bool) {
	var (
		post *models.Post
		err error
	)
	id := c.Param("id")
	postID, err := strconv.ParseUint(id,10,64)
	if err != nil {
		msg := fmt.Sprintf("parse uint error:%v",err)
		Logger.Fatal(msg)
	}
	if isPublish {
		post,err = models.GetPostByIDAndPublished(postID,isPublish)
	} else {
		post,err = models.GetPostByID(postID)
	}

	if err != nil {
		c.AbortWithError(http.StatusNotFound,err)
		return
	}
	tags ,err := models.ListTagByPostID(post.ID)
	if err != nil {
		msg := fmt.Sprintf("list tag by postID error:%v",err)
		Logger.Fatal(msg)
	}
	post.Tags = tags
	content := post.Content
	comments, err := models.ListCommentsByPostID(int(postID))
	if err != nil {
		msg := fmt.Sprintf("list comments by postID error:%v",err)
		Logger.Fatal(msg)
	}
	gitHubUser, _ := c.Get(models.CONTEXT_GIT_USER_KEY)
	policy := bluemonday.UGCPolicy()
	render := blackfriday.HtmlRenderer(commonHtmlFlags,"","")
	unsafe := blackfriday.Markdown([]byte(content),render,commonExtensions)
	contentHtml:=template.HTML(policy.Sanitize(string(unsafe)))

	var pages int
	if len(comments) % 10 == 0 {
		pages = len(comments) / 10
	} else {
		pages =  len(comments) / 10 + 1
	}
	hh := utils.HH{
		Post:post,
		Comments:comments,
		Githubuser:gitHubUser,
		Pages:pages,
		CommentNum:len(comments),
	}
	commentsHTML,_ := utils.RenderAllComment(hh)
	res := template.HTML(commentsHTML)

	relatePosts := GetPosts(int64(postID))


	c.HTML(http.StatusOK,"front/post.html",gin.H{
		"Post":post,
		"contentHtml":contentHtml,
		"Comments": comments,
		"Githubuser":gitHubUser,
		"Pages": pages,
		"CommentNum":len(comments),
		"commentsHTML":res,
		"relatePosts": relatePosts,
	})
}

func GetPosts(postID int64) []*models.Post {
	tags, err := models.ListTagByPostID(postID)
	if err != nil {
		msg := fmt.Sprintf("list tag by postID error:%v",err)
		Logger.Fatal(msg)
	}
	var tagids []int64
	for _,tag := range tags {
		tagids = append(tagids,int64(tag.ID))
	}
	posts := models.GetPostsByTags(postID,tagids)
	result := utils.RandomGetArray(posts,4)
	return result
}
