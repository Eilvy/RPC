package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go_code/RPC/utils"
	"gorm.io/gorm"
	"time"
)

func CreateUser(c *gin.Context) {
	var user utils.User
	var test utils.User

	err := c.ShouldBind(&user)
	if err != nil {
		utils.Failed(c)
		utils.Logger.Println("获取输入用户数据失败:", err.Error()) //写入错误日志
		return
	}
	if user.Password == "" {
		c.JSON(200, gin.H{"msg": "请必须设置密码"})
		return
	}
	//err = utils.DB.QueryRow("select username from users where username=?", user.Username).Scan(&test.Username)
	result := utils.DB.Where("username=?", user.Username).First(&test)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{"msg": "数据库中查无此人，请重新输入"})
			return
		} else {
			fmt.Println("查询出错，err:", result.Error.Error())
			return
		}
	}
	//if err == nil {
	//出现查重数据
	//	c.JSON(200, gin.H{"msg": "用户名重复，请重新输入"})
	//	utils.Failed(c)
	//	return
	//} else if !errors.Is(err, sql.ErrNoRows) {
	//查重过程出错，程序出错
	//	c.JSON(200, gin.H{"msg": "database error", "err": err.Error()})
	//	utils.Failed(c)
	//	return
	//}
	//_, err = utils.DB.Exec("insert into users (username,password,LongURL,ShortURL,Prise) values (?,?,?,?,?)", user.Username, user.Password, user.LongURL, user.ShortURL, user.Prise)
	result = utils.DB.Create(&user)
	if result.Error != nil {
		c.JSON(200, gin.H{"mag": "传入数据库数据失败", "err": result.Error.Error()})
		return
	}
	//if err != nil {
	//	utils.Failed(c)
	//	utils.Logger.Println("写入用户数据失败:", err.Error())
	//	return
	//}
	c.JSON(200, gin.H{"msg": "success"})
}

func GetToken(c *gin.Context) {
	//var Password string
	//var test string
	username := c.Query("username")
	password := c.Query("password")
	//从数据库中查找此人是否存在
	//err := utils.DB.QueryRow("select username from users where (username=?) and password=?", username, password).Scan(&test)
	result := utils.DB.Where("username=?", username).First(&utils.Test)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{"msg": "数据库中查无此人，请重新输入"})
			return
		} else {
			fmt.Println("查询出错，err:", result.Error.Error())
			return
		}
	}

	//if err != nil {
	//	fmt.Printf("error:%s\n", err.Error())
	//	c.JSON(200, gin.H{"msg": "数据库中查无此人，请重新输入"})
	//	utils.Failed(c)
	//	return
	//}

	//err := utils.DB.QueryRow("select password from users where username=?", test).Scan(&Password)
	result = utils.DB.Where("username=?", username).First(&utils.Test)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{"msg": "数据库中查无此人，请重新输入"})
			return
		} else {
			fmt.Println("查询出错，err:", result.Error.Error())
			return
		}
	}
	//if err != nil {
	//	utils.Failed(c)
	//	utils.Logger.Println("密码读取错误:", err.Error())
	//	c.JSON(200, gin.H{"data": "密码数据读取错误", "err": err.Error()})
	//	return
	//}
	if utils.Test.Password != password {
		utils.Failed(c)
		c.JSON(200, gin.H{"data": "密码输入错误"})
		c.Abort()
		return
	}
	//token claim声明
	claims1 := jwt.MapClaims{
		"username": username,
		"password": password,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}
	claims2 := jwt.MapClaims{
		"username": username,
		"password": password,
		"exp":      time.Now().Add(14 * 24 * time.Hour).Unix(), //一个14天的refresh_token
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims1)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims2)
	signedToken, err := token.SignedString(utils.Key)
	if err != nil {
		c.JSON(200, gin.H{"msg": "生成token失败"})
		utils.Failed(c)
		return
	}
	signedRefreshToken, err := refreshToken.SignedString(utils.Key)
	if err != nil {
		c.JSON(200, gin.H{"msg": "生成refreshToken失败"})
		utils.Failed(c)
		return
	}
	c.Request.Header.Set("Authorization", "Bearer "+signedToken) //设置请求头（无法在postman里查询到）
	fmt.Println(c.Request.Header.Get("Authorization"))
	fmt.Printf("登录用户为：%s\n", username)
	c.JSON(200, gin.H{"msg": "success", "data": gin.H{"refreshToken": signedRefreshToken, "Token": signedToken}})
}

func RefreshToken(c *gin.Context) {
	//var test string
	username := c.MustGet("username").(string)
	//err := utils.DB.QueryRow("select username from users where username=?", username).Scan(&test)
	//if err != nil {
	//	c.JSON(200, gin.H{"msg": "数据库中查无此人，请重新输入"})
	//	utils.Logger.Println("err:", err.Error())
	//	utils.Failed(c)
	//	return
	//}
	result := utils.DB.Where("username=?", username).First(&utils.Test)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{"mas": "查无此人"})
			return
		} else {
			c.JSON(200, gin.H{"msg": "查询数据库失败", "err": result.Error.Error()})
			return
		}
	}
	claims1 := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}
	claims2 := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(14 * 24 * time.Hour).Unix(), //一个14天的refresh_token
	}
	nToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims1)
	nRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims2)
	signedToken, err := nToken.SignedString(utils.Key)
	if err != nil {
		c.JSON(200, gin.H{"msg": "生成token失败"})
		utils.Failed(c)
		return
	}
	signedRefreshToken, err := nRefreshToken.SignedString(utils.Key)
	if err != nil {
		c.JSON(200, gin.H{"msg": "生成refreshToken失败"})
		utils.Logger.Println("生成refreshToken失败", err.Error())
		utils.Failed(c)
		return
	}
	c.Request.Header.Set("Authorization", "Bearer "+signedToken)
	fmt.Printf("登录用户为：%s\n", username)
	c.JSON(200, gin.H{"msg": "success", "data": gin.H{"refreshToken": signedRefreshToken, "Token": signedToken}})
}

func AddLongURL(c *gin.Context) {
	longURL := c.PostForm("longURL")
	username := c.PostForm("username")
	exists, err := utils.Redis.Exists(context.Background(), username).Result()
	if err != nil {
		utils.Logger.Println("读取redis失败,err:", err.Error())
		return
	}

	if exists == 1 {
		err := utils.Redis.Del(context.Background(), username).Err()
		if err != nil {
			utils.Logger.Println("删除缓存失败,err:", err.Error())
			return
		}
	}
	//存入数据库mysql
	//var test string
	//err = utils.DB.QueryRow("select username from users where username=?", username).Scan(&test)
	//if err != nil {
	//	utils.Logger.Println("查询数据库用户名失败,err:", err.Error())
	//	return
	//}
	result := utils.DB.Where("username=?", username).First(&utils.Test)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{"mas": "查无此人"})
			return
		} else {
			c.JSON(200, gin.H{"msg": "查询数据库失败", "err": result.Error.Error()})
			return
		}
	}
	//_, err = utils.DB.Exec("update users set LongURL=? where username=?", longURL, username)
	//if err != nil {
	//	utils.Logger.Println("写入LongURL失败,err:", err.Error())
	//	return
	//}

	result = utils.DB.Model(&utils.Test).Where("username=?", username).Update("long_url", longURL)
	if result.Error != nil {
		c.JSON(200, gin.H{"msg": "写入数据库失败", "err": result.Error.Error()})
		return
	}

	//存入缓存redis
	err = utils.Redis.Set(context.Background(), username, longURL, 0).Err()
	if err != nil {
		utils.Logger.Println("长链接存入缓存失败,err:", err.Error())
		return
	}
	c.JSON(200, gin.H{"msg": "success"})
}
