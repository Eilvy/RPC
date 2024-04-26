package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	item "go_code/RPC/kitex_gen/example/shop/item"
	"go_code/RPC/utils"
	"time"
)

// ItemServiceImpl implements the last service interface defined in the IDL.
type ItemServiceImpl struct{}

// GetItem implements the ItemServiceImpl interface.
func (s *ItemServiceImpl) GetItem(ctx context.Context, req *item.GetItemReq) (resp *item.GetItemResp, err error) {
	// TODO: Your code here...
	redisCli := redis.NewClient(&redis.Options{
		Addr:        "redis-14520.c299.asia-northeast1-1.gce.cloud.redislabs.com:14520",
		Password:    "rPYdtUeiD5CeJSqcGZdoyHDd6Ou2uApa",
		DB:          0,
		DialTimeout: time.Second * 5,
	})
	resp = item.NewGetItemResp()
	resp.Item = item.NewItem()
	resp.Item.Title = "Kitex"
	resp.Item.Description = "Kitex is an excellent framework!"
	resp.ShortURL, err = utils.Shorten(redisCli, req.LongURL, 100)
	if err != nil {
		utils.Logger.Println("转换短链错误,err:", err.Error())
		return
	}
	return
}
