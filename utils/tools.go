package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	username := c.MustGet("username")

	longUrl, err := Redis.Get(context.Background(), fmt.Sprintf(RedisShortUrl, sid)).Result()
	if err != nil {
		Logger.Println("重定向查找长链失败,err:", err.Error())
		c.String(http.StatusNotFound, "short URL is out of the expireTime or undefined")

		//从redis里查找长链失败的同时删除mysql的短链
		result := DB.Model(&Test).Where("username=?", username).Update("short_url", "")
		if result.Error != nil {
			c.JSON(200, gin.H{"msg": "删除用户短链错误", "err": result.Error.Error()})
			Logger.Println("删除用户短链错误,err:", result.Error.Error())
			return
		}
		return
	}

	//查询数据库
	result := DB.Where("short_url=?", sid).First(&Test)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{"msg": "用户短链不存在"})
			return
		} else {
			c.JSON(200, gin.H{"msg": "查询出错", "err": result.Error.Error()})
			Logger.Println("查询出错，err:", result.Error.Error())
			return
		}
	}
	//访问数加一
	Test.Prise += 1
	result = DB.Where("short_url=?", sid).Update("prise", Test.Prise)
	if result.Error != nil {
		c.JSON(200, gin.H{"msg": "写入访问数错误", "err": result.Error.Error()})
		Logger.Println("写入访问数错误,err:", result.Error.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, longUrl)

} //短链重定向工具函数
