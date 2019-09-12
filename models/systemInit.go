package models

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

const (
	SESSION_KEY = "UserID"
	CONTEXT_USER_KEY = "User"
	SESSION_GITHUB_STATE = "GITHUB_STATE" // github state session key
)

var (
	DB *gorm.DB
	RedisPool *redis.Pool
	Conf = new(Config)
)

type Config struct {
	RunMode string
	General struct{
		Addr string
		DSN string
		RedisUrl string
		SessionSecret string
		LogOutEnabled bool
	}
	GitHub struct{
		ClientID string
		ClientSecret string
		AuthUrl string
		RedirectUrl string
		TokenUrl string
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func InitDB()(err error) {
	db,err := gorm.Open("mysql",Conf.General.DSN)
	if err == nil {
		log.Println("connect db success")
		DB = db
		//DB.LogMode(true)
		DB.AutoMigrate(&Comment{}, &Post{},&PostTag{}, ReactItem{},Tag{},User{})
	}
	return
}

//redis存文章的内容格式例子：posts/1/props/content

func initRedis(server, password string)*redis.Pool {
	return &redis.Pool{
		MaxIdle:64,
		MaxActive:100,
		IdleTimeout:240 * time.Second,
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", server)
			if err != nil {
				return nil,err
			}
			/*
				if _, err = conn.Do("AUTH", password);err != nil {
					conn.Close()
					return nil,err
				}
			*/
			return conn, err
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}


func init() {
	data, err := ioutil.ReadFile("config/config.yaml")
	checkError(err)
	err = yaml.Unmarshal(data, Conf)
	checkError(err)
	log.Println(Conf)
	err = InitDB()
	checkError(err)
	RedisPool = initRedis(Conf.General.RedisUrl,"")
	fmt.Println(RedisPool)

}