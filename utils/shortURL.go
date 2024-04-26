package utils

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mattheath/base62"
	"github.com/redis/go-redis/v9"
	"time"
)

func CreateShortURL(c *gin.Context) {
	username := c.MustGet("username").(string) //从中间件（token）获取用户名
	exists, err := Redis.Exists(context.Background(), username).Result()
	if err != nil {
		Logger.Println("查询redis失败,err:", err.Error())
		return
	}
	if exists == 1 {
		longURL, err := Redis.Get(context.Background(), username).Result()
		if err != nil {
			Logger.Println("查询用户长链失败,err:", err.Error())
			return
		}
		sid, err := Shorten(Redis, longURL, 100)
		if err != nil {
			Logger.Println("长链缩短失败,err:", err.Error())
			return
		}
		result := DB.Where("username=?", username).First(&Test)
		if result.Error != nil {
			Logger.Println("读取数据库错误,err:", result.Error.Error())
			c.JSON(200, gin.H{"msg": "读取数据库错误", "err": result.Error.Error()})
			return
		}

		//将短链sid写入数据库
		result = DB.Model(&Test).Where("username=?", username).Update("short_url", sid)
		if result.Error != nil {
			Logger.Println("写入用户短链错误,err:", result.Error.Error())
			c.JSON(200, gin.H{"msg": "写入用户短链错误", "err": result.Error.Error()})
			return
		}

		c.JSON(200, gin.H{"msg": "success", "shortURL": sid, "longURL": longURL, "user": username})
		return
	}
	c.JSON(200, gin.H{"msg": "用户长链不存在"})
}

func Shorten(RedisCli *redis.Client, url string, expireTime int64) (string, error) {
	//传入url为长连接，
	id, err := RedisCli.Incr(context.Background(), RedisKeyUrlID).Result() //从redis获取自增ID
	if err != nil {
		Logger.Println("获取自增ID失败,err:", err.Error())
		return "", err
	}
	sid := base62.EncodeInt64(Offset + id) //用base62的方式进行转换//因为base64里包含+，/两个特殊符号对URL不友好不选择//用哈希缩的不会很短

	err = RedisCli.Set(context.Background(), fmt.Sprintf(RedisShortUrl, sid), url, time.Second*time.Duration(expireTime)).Err()
	//fmt.Sprintf(RedisShortUrl, sid)将sid写入RedisShortUrl中的%s部分并作为redis的一个键
	//url是上面键的值，将sid和长链url用键值对的形式绑定
	//time.Second*time.Duration(expireTime)将int64类型的expireTime转化为时间段并转成以秒为单位
	if err != nil {
		Logger.Println("短链绑定长链失败,err:", err.Error())
		return "", err
	}
	return "http://127.0.0.1/" + sid, nil
}
