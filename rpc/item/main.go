package main

import (
	item "go_code/RPC/kitex_gen/example/shop/item/itemservice"
	"log"
)

func main() {
	svr := item.NewServer(new(ItemServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}

}
