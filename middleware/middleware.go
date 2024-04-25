package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go_code/RPC/utils"
	"gorm.io/gorm"
	"strings"
)

func JWTAuthorMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			c.JSON(200, gin.H{
				"status": 20000,
				"info":   "failed",
				"data":   "auth为空或者token已过期",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authorization, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(200, gin.H{
				"status": 20000,
				"info":   "failed",
				"data":   "请求头中auth格式有误或者token已过期",
			})
			c.Abort()
			return
		}
		mc, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Failed(c)
			c.JSON(200, gin.H{"data": err.Error()})
			c.Abort()
			return
		}
		c.Set("username", mc.Username)
		//var user_id string
		//err = utils.DB.QueryRow("select users.id from users where username=?", mc.Username).Scan(&user_id)
		//if err != nil {
		//	utils.Failed(c)
		//	c.JSON(200, gin.H{"data": "获取用户id失败", "err": err.Error()})
		//	c.Abort()
		//
		//}

		result := utils.DB.Where("username=?", mc.Username).First(&utils.Test)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(200, gin.H{"msg": "数据库查无此人"})
				return
			} else {
				c.JSON(200, gin.H{"msg": "查询失败", "err": result.Error.Error()})
				return
			}
		}
		c.Set("user_id", utils.Test.Id)
		c.Next()
	}
} //用户中间件
