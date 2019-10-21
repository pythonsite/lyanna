package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"lyanna/models"
	"lyanna/utils"
)


//var Logger = models.Logger

func GetRss(c *gin.Context) {
	now := utils.GetCurrentTime()
	feed := &feeds.Feed{
		Title:"Fan's Blog",
		Link:&feeds.Link{Href:"www.syncd.com"},
		Description:"Fan's Blog about golang,python",
		Author:&feeds.Author{Name:"fan",Email:"hjzhaofan@163.com"},
		Created:now,
	}
	feed.Items = make([]*feeds.Item, 0)
	posts, err := models.ListPublishedPost("")
	if err != nil {
		msg := fmt.Sprintf("list published posts err:%v",err)
		Logger.Fatal(msg)
	}
	for _, post := range posts {
		item := &feeds.Item{
			Id:          fmt.Sprintf("%s/post/%d", "www.syncd.com", post.ID),
			Title:       post.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/post/%d", "www.syncd.com", post.ID)},
			Description: post.Summary,
			Created:     now,
		}
		feed.Items = append(feed.Items, item)
	}
	rss, err := feed.ToRss()
	if err != nil {
		msg := fmt.Sprintf("feed to rss err:%v",err)
		Logger.Fatal(msg)
	}
	c.Writer.WriteString(rss)

}