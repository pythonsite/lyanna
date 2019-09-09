package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"lyanna/models"
	"lyanna/utils"
	"net/http"
	"strconv"
)

func UserLogin(c *gin.Context) {
	var (
		err error
		user *models.User
	)
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.HTML(http.StatusOK,"admin/login.html", gin.H{
			"msg":"username or password not null",
		})
		return
	}
	user, err = models.GetUserByName(username)
	if err != nil || user.PassWord != utils.Md5(username + password) {
		c.HTML(http.StatusOK, "admin/login.html", gin.H{
			"msg": "invalid username or passwrod",
		})
		return
	}
	if !user.Active {
		c.HTML(http.StatusOK,"admin/login.html",gin.H{
			"msg":"user is not active",
		})
		return
	}
	s := sessions.Default(c)
	s.Clear()
	s.Set(models.SESSION_KEY,user.ID)
	s.Save()
	c.Redirect(http.StatusMovedPermanently,"/admin/index")
}

func UserList(c *gin.Context) {
	users, _ := models.ListUsers()
	user, _ := c.Get(models.CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/list_user.html",gin.H{
		"users": users,
		"user": user,
		"user_count":len(users),
	})
}

func PostUserEdit(c *gin.Context) {
	users, _ := models.ListUsers()
	id := c.Param("id")
	uID, err := strconv.ParseUint(id,10,64)
	if err != nil {
		return
	}
	name := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	active :="on" == c.PostForm("active")

	user := &models.User{
		Name:name,
		Email:email,
		PassWord:password,
		Active:active,
	}
	user.ID = uID
	err = user.Update()
	if err != nil {
		c.HTML(http.StatusOK, "admin/list_user.html",gin.H{
			"user":user,
			"msg": err.Error(),
		})
	}
	c.HTML(http.StatusOK,"admin/list_user.html",gin.H{
		"users":users,
		"user":user,
		"user_count":len(users),
		"msg":"User was successfully updated.",
	})
}

func GetEditUser(c *gin.Context) {
	id := c.Param("id")
	uID, err := strconv.ParseUint(id,10,64)
	if err != nil {
		return
	}
	user, _ := models.GetUserByID(uID)
	log.Println(user)
	if user == nil {
		c.Redirect(http.StatusMovedPermanently,"/admin/users")
	}
	c.HTML(http.StatusOK, "admin/user.html",gin.H{
		"user":user,
	})
}

func GetCreateUser(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/user.html",nil)
}

func PostCreateUser(c *gin.Context) {
	name := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	active :="on" == c.PostForm("active")
	md5Password := utils.Md5(name + password)
	user := &models.User{
		Name:name,
		Email:email,
		PassWord:md5Password,
		Active:active,
	}
	err := user.Insert()
	if err != nil {
		return
	}
	users, _ := models.ListUsers()
	c.HTML(http.StatusOK,"admin/list_user.html",gin.H{
		"users":users,
		"user":user,
		"user_count":len(users),
		"msg":"User was successfully created.",
	})
}

