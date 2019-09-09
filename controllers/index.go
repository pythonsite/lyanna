package controllers

import (
	"github.com/gin-gonic/gin"
	"lyanna/models"
	"net/http"
)

func IndexGet(c *gin.Context) {
	c.HTML(http.StatusOK,"index.html",gin.H{
		"title":"hello golang",
	})
}

func AdminLogin(c *gin.Context) {
	c.HTML(http.StatusOK,"admin/login.html",nil)
}

func AdminIndex(c *gin.Context) {
	user, _ := c.Get(models.CONTEXT_USER_KEY)
	if user == nil {
		c.Redirect(http.StatusMovedPermanently,"/admin/login")
	}
	c.HTML(http.StatusOK, "admin/index.html",nil)
}
