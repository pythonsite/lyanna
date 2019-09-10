package models

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func getKey(postID int) string {
	return fmt.Sprintf(RedisPostKey, postID)
}

func SetContent(postID int, value string) {
	key := getKey(postID)
	conn := RedisPool.Get()
	_,_ = conn.Do("set",key,value)
}

func GetContent(postID int) string {
	key := getKey(postID)
	conn := RedisPool.Get()
	value, _ := redis.String(conn.Do("get",key))
	return value
}