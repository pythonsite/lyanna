package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Index(c *gin.Context) {
	log.Println("aaaa")
	c.HTML(http.StatusOK,"front/index.html",nil)
}
