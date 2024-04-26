package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/gin-gonic/gin"
	"go_code/RPC/config"
	"go_code/RPC/kitex_gen/example/shop/item"
	"go_code/RPC/kitex_gen/example/shop/item/itemservice"
	"go_code/RPC/middleware"
	"go_code/RPC/utils"
	"log"
	"os"
	"time"
)

func main() {
	//用golang 原生log库启动错误日志
	folderPath2 := utils.LogAdd()
	logFile, err := os.Create(folderPath2) //在本项目目录下创建文件夹loggerFold并在其中创建日志文件error.log
	if err != nil {
		fmt.Println("add logFile error:", err.Error())
		return
	}
	defer logFile.Close()
	utils.Logger = log.New(logFile, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile)

	config.RedisConnect()
	config.MysqlConnect()
	//defer utils.DB.Close()

	utils.Router = gin.Default()
	utils.Router.POST("/createUser", config.CreateUser)                                         //创建新用户
	utils.Router.POST("/getToken", config.GetToken)                                             //获取用户Token
	utils.Router.GET("/refreshToken", middleware.JWTAuthorMiddleware(), config.RefreshToken)    //刷新用户Token
	utils.Router.POST("/longURLAdd", middleware.JWTAuthorMiddleware(), config.AddLongURL)       //注册用户写入longURL
	utils.Router.GET("/createShortURL", middleware.JWTAuthorMiddleware(), utils.CreateShortURL) //获取用户对应longURL的shortURL
	utils.Router.GET("/:sid", middleware.JWTAuthorMiddleware(), utils.RedirectHandler)          //对短链进行重定向

	c, err := itemservice.NewClient("example.shop.item.exe", client.WithHostPorts("0.0.0.0:8888"))
	if err != nil {
		utils.Logger.Println("启动kitex服务error:", err.Error())
		return
	}
	utils.Cli = c

	utils.Router.GET("/api/item", middleware.JWTAuthorMiddleware(), Handler)

	if err := utils.Router.Run(":80"); err != nil {
		utils.Logger.Println("启动gin服务error:", err.Error())
		return
	}
}

func Handler(c *gin.Context) {
	req := item.NewGetItemReq()
	username := c.MustGet("username").(string) //从中间件（token）获取用户名
	exists, err := utils.Redis.Exists(context.Background(), username).Result()
	if err != nil {
		utils.Logger.Println("查询redis失败,err:", err.Error())
		return
	}
	if exists == 1 {
		req.LongURL, err = utils.Redis.Get(context.Background(), username).Result()
		if err != nil {
			utils.Logger.Println("查询用户长链失败,err:", err.Error())
			return
		}
		//sid, err := utils.Shorten(longURL, 100)
		resp, err := utils.Cli.GetItem(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
		if err != nil {
			utils.Logger.Println("调用kitex失败,err:", err.Error())
			log.Fatal(err)
			return
		}
		c.String(200, resp.String())
		if err != nil {
			utils.Logger.Println("长链缩短失败,err:", err.Error())
			return
		}
		//c.JSON(200, gin.H{"msg": "success", "shortURL": utils.sid, "longURL": longURL, "user": username})
		return
	}
	// 调用kitex客户端
	//resp, err := utils.Cli.GetItem(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//c.String(200, resp.String())
}
