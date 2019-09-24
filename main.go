package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"lyanna/controllers"
	"lyanna/models"
	"lyanna/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	gin.SetMode(models.Conf.RunMode)
	router := gin.Default()
	setTemplate(router)
	setSessions(router)
	router.Use(ShareData())
	router.Static("/static", filepath.Join(getCurrentDirectory(), "./static"))

	router.GET("/", controllers.Index)
	router.GET("/tags",controllers.Tags)
	router.GET("/tag/:id",controllers.Tag)
	router.GET("/post/:id",controllers.GetPost)

	router.GET("/archives",controllers.Archives)
	router.GET("/archives/:year",controllers.ArchivesByYear)

	router.GET("/oauth2/auth",controllers.AuthGet)
	router.GET("/oauth2",controllers.Oauth2Callback)
	router.GET("/oauth2/auth/post/:id", controllers.AuthGet)
	router.GET("/admin/login",controllers.AdminLogin)
	router.POST("/admin/login",controllers.UserLogin)
	router.POST("/api/publish/:id", controllers.PostPublish)
	router.DELETE("/api/publish/:id",controllers.DeletePublish)
	router.POST("/comment/post/:id", controllers.CreateComment)
	router.POST("/comment/markdown",controllers.CommentMarkdown)
	router.GET("/comments/post/:id",controllers.Comments)
	router.GET("/rss",controllers.GetRss)
	router.GET("/page/:aboutme",controllers.AboutMe)
	router.GET("/search",controllers.GetSearch)
	router.GET("/json/search", controllers.PostSearch)
	router.GET("/pages/:page",controllers.PostPage)

	admin := router.Group("/admin")
	admin.Use(AdminRequired())
	{
		admin.GET("/posts", controllers.PostIndex)

		admin.GET("/post/edit/:id", controllers.GetEditPost)
		admin.POST("/post/edit/:id",controllers.UpdatePost)

		admin.GET("/post/new",controllers.GetNewPost)
		admin.POST("/post/new", controllers.AddPost)

		admin.GET("/posts/page/:page", controllers.AdminPostPage)
		admin.GET("/users/page/:page", controllers.AdminUserPage)

		admin.GET("/post/preview/:id", controllers.PreviewGetPost)
		admin.GET("/",controllers.AdminIndex)

		admin.GET("/users",controllers.UserList)
		admin.GET("/user/edit/:id",controllers.GetEditUser)
		admin.POST("/user/edit/:id",controllers.PostUserEdit)
		admin.GET("/user/new",controllers.GetCreateUser)
		admin.POST("/user/new", controllers.PostCreateUser)
	}

	err := router.Run(models.Conf.General.Addr)
	if err != nil {
		log.Fatal(err)
	}
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir,"\\", "/", -1)
}

func setTemplate(engine *gin.Engine) {
	funcMap := template.FuncMap{
		"dateFormat": utils.DateFormat,
		"genList":utils.GenList,
		"add": utils.Add,
		"GetMapValue":utils.GetMapValue,
	}
	engine.SetFuncMap(funcMap)
	engine.LoadHTMLGlob(filepath.Join(getCurrentDirectory(), "./views/**/*"))
}

func setSessions(router *gin.Engine) {
	store := cookie.NewStore([]byte(models.Conf.General.SessionSecret))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 2 * 86400, Path: "/"})
	router.Use(sessions.Sessions("gin-session",store))
}

func ShareData() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if uID := session.Get(models.SESSION_KEY);uID != nil {
			user,err := models.GetUserByID(uID)
			if err == nil {
				c.Set(models.CONTEXT_USER_KEY, user)
			}
			gitUser, err := models.GetGitUserByGid(uID)
			if err == nil {
				c.Set(models.CONTEXT_GIT_USER_KEY, gitUser)
			}
			if models.Conf.General.LogOutEnabled {
				c.Set("LogOutEnabled", true)
			}
			c.Next()
		}
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user,_ := c.Get(models.CONTEXT_USER_KEY);user != nil {
			if _, ok := user.(*models.User);ok {
				c.Next()
				return
			}
		}
		c.HTML(http.StatusForbidden,"errors/error.html",gin.H{
			"message":"Forbidden!",
		})
		c.Abort()
	}
}

