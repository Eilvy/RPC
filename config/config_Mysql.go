package config

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go_code/RPC/utils"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MysqlConnect() {
	username := "root"
	password := "root"
	host := "127.0.0.1"
	port := 3306
	Dbname := "rpcsystem"
	timeout := "10s"
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)

	db, err := gorm.Open(mysql.Open(dns))
	if err != nil {
		utils.Logger.Println("连接数据库mysql失败，err:", err.Error())
		fmt.Println("连接数据库mysql失败，err:", err.Error())
		return
	}
	utils.DB = db
}
