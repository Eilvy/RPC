package utils

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go_code/RPC/kitex_gen/example/shop/item/itemservice"
	"log"
)

var (
	Redis  *redis.Client
	DB     *sql.DB
	Key    = []byte("leiyv")
	Logger *log.Logger
	Cli    itemservice.Client
	Router *gin.Engine
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

type Users struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Prise    int    `form:"prise"`
	LongURL  string `form:"longUrl"`
	ShortURL string `form:"shortUrl"`
}

type Shortage interface {
	Shorten(url string, expireTime int64) (string, error) //长链接转短链
	unShorten(sid string) (string, error)                 //短链返回成长链
}
