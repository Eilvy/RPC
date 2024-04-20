package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func Failed(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "failed"})
}

func ParseToken(tokenString string) (*MyClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return Key, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
} // 解析token，用于中间件解析

func LogAdd() (folderPath2 string) {
	currentDir, err := os.Getwd() //获取当前项目目录
	if err != nil {
		fmt.Println("get current dir failed:", err.Error())
		return
	}

	loggerFold := "loggerFold"

	folderPath1 := filepath.Join(currentDir, loggerFold)

	_, err = os.Stat(folderPath1)
	if os.IsNotExist(err) { //创建文件夹失败
		Err := os.Mkdir(loggerFold, 0755)
		if Err != nil {
			fmt.Println("add loggerFold error:", Err.Error())
			return
		}
	} else if err != nil {
		fmt.Println("find folder error:", err.Error())
		return
	}
	//创建文件夹成功
	folderPath2 = filepath.Join(folderPath1, "error.log")
	return folderPath2
} //简易的错误日志，用golang原生log库实现//此处为创建文件夹

func RedirectHandler(c *gin.Context) {
	sid := c.Param("sid")

	longUrl, err := Redis.Get(context.Background(), fmt.Sprintf(RedisShortUrl, sid)).Result()
	if err != nil {
		Logger.Println("重定向查找长链失败,err:", err.Error())
		c.String(http.StatusNotFound, "short URL is out of the expireTime or undefined")
		return
	}
	c.Redirect(http.StatusSeeOther, longUrl)

} //短链重定向工具函数
