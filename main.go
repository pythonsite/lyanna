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

	router.GET("/admin/login",controllers.AdminLogin)
	router.POST("/admin/login",controllers.UserLogin)
	router.GET("/admin",controllers.AdminIndex)
	router.GET("/admin/users",controllers.UserList)
	router.GET("/admin/edit/:id",controllers.GetEditUser)
	router.POST("/admin/edit/:id",controllers.PostUserEdit)
	router.GET("/admin/user/new",controllers.GetCreateUser)
	router.POST("/admin/user/new", controllers.PostCreateUser)

	router.GET("/admin/posts", controllers.PostIndex)

	router.POST("/api/publish/:id", controllers.PostPublish)
	router.DELETE("/api/publish/:id",controllers.DeletePublish)

	router.GET("/admin/post/edit/:id", controllers.GetEditPost)
	router.GET("/admin/post/new",controllers.GetNewPost)
	router.POST("/admin/post/edit/:id",controllers.PostCreatePost)

	//admin := router.Group("/admin")
	//{
	//
	//}


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
			if models.Conf.General.LogOutEnabled {
				c.Set("LogOutEnabled", true)
			}
			c.Next()
		}
	}
}

