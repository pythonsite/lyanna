package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

const (
	SESSION_KEY = "UserID"
	CONTEXT_USER_KEY = "User"
)

var (
	DB *gorm.DB
	Conf = new(Config)
)

type Config struct {
	RunMode string
	General struct{
		Addr string
		DSN string
		SessionSecret string
		LogOutEnabled bool
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
		DB.LogMode(true)
		DB.AutoMigrate(&Comment{}, &Post{},&PostTag{}, ReactItem{},Tag{},User{})
	}
	return
}


func init() {
	data, err := ioutil.ReadFile("config/config.yaml")
	checkError(err)
	err = yaml.Unmarshal(data, Conf)
	checkError(err)
	log.Println(Conf)
	err = InitDB()
	checkError(err)
}