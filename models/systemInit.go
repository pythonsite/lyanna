package models

import (
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

const (
	SESSION_KEY = "UserID"
	CONTEXT_USER_KEY = "User"
	CONTEXT_GIT_USER_KEY = "GitUser"
	SESSION_GITHUB_STATE = "GITHUB_STATE" // github state session key
)

var (
	DB *gorm.DB
	RedisPool *redis.Pool
	Conf = new(Config)
	Logger *zap.Logger
)

type Config struct {
	RunMode string
	General struct{
		Addr string
		DSN string
		SessionSecret string
		LogOutEnabled bool
		PerPage int
	}
	GitHub struct{
		ClientID string
		ClientSecret string
		AuthUrl string
		RedirectUrl string
		TokenUrl string
	}
	Log struct{
		LogPath string
		MaxSize int
		MaxAge int
		Compress bool
		MaxBackups int
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
		DB.AutoMigrate(&Comment{}, &Post{},&PostTag{}, &ReactItem{}, &Tag{}, &User{}, &GitHubUser{})
	}
	return
}

//redis存文章的内容格式例子：posts/1/props/content

//func initRedis(server, password string)*redis.Pool {
//	return &redis.Pool{
//		MaxIdle:64,
//		MaxActive:100,
//		IdleTimeout:240 * time.Second,
//		Dial: func() (conn redis.Conn, err error) {
//			conn, err = redis.Dial("tcp", server)
//			if err != nil {
//				return nil,err
//			}
//			/*
//				if _, err = conn.Do("AUTH", password);err != nil {
//					conn.Close()
//					return nil,err
//				}
//			*/
//			return conn, err
//		},
//		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
//			if time.Since(t) < time.Minute {
//				return nil
//			}
//			_, err := conn.Do("PING")
//			return err
//		},
//	}
//}

// 初始化日志配置
func initLog() {
	hook := lumberjack.Logger{
		Filename:Conf.Log.LogPath,  //日志文件路径
		MaxSize: Conf.Log.MaxSize,// 每个日志的大小，单位是M
		MaxAge: Conf.Log.MaxAge, // 文件被保存的天数
		Compress:Conf.Log.Compress, // 是否压缩
		MaxBackups:Conf.Log.MaxBackups, // 保存多少个文件备份
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:"Time",
		LevelKey:"Level",
		NameKey: "Logger",
		CallerKey: "Caller",
		MessageKey:"Msg",
		LineEnding: zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime: zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout),zapcore.AddSync(&hook)),
		atomicLevel,
	)
	caller := zap.AddCaller()
	development := zap.Development()
	filed := zap.Fields(zap.String("service","blog"))
	Logger = zap.New(core,caller, development,filed)
}

func init() {
	data, err := ioutil.ReadFile("config/config.yaml")
	checkError(err)
	err = yaml.Unmarshal(data, Conf)
	checkError(err)
	initLog()
	Logger.Info("init config and log config success")
	err = InitDB()
	checkError(err)
	//RedisPool = initRedis(Conf.General.RedisUrl,"")
	//fmt.Println(RedisPool)

}