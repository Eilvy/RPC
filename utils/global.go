package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go_code/RPC/kitex_gen/example/shop/item/itemservice"
	"gorm.io/gorm"
	"log"
)

var (
	Redis  *redis.Client
	DB     *gorm.DB
	Key    = []byte("leiyv")
	Logger *log.Logger
	Cli    itemservice.Client
	Router *gin.Engine
	Test   User
)

const (
	Offset        = 1000000 //初始序号
	RedisKeyUrlID = "url:global:id"
	RedisShortUrl = "url:short:%s"
)

type MyClaims struct {
	Username string `json:"Username"`
	jwt.StandardClaims
}

//	type Users struct {
//		Username string `form:"username"`
//		Password string `form:"password"`
//		Prise    int    `form:"prise"`
//		LongURL  string `form:"longUrl"`
//		ShortURL string `form:"shortUrl"`
//	}
type User struct {
	Id       uint   `gorm:"primaryKey;autoIncrement"` //自增的主键ID
	Username string `gorm:"type:char;size:100"`
	Password string `gorm:"type:char;size:255"`
	Prise    int    `gorm:"type:int;"`
	LongURL  string `gorm:"type:char;size:255"`
	ShortURL string `gorm:"type:char;size:200"`
}

type Shortage interface {
	Shorten(url string, expireTime int64) (string, error) //长链接转短链
	unShorten(sid string) (string, error)                 //短链返回成长链
}
