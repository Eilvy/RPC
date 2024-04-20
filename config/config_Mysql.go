package config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go_code/RPC/utils"
)

func MysqlConnect() {
	utils.DB, _ = sql.Open("mysql", "root:lei_yv_0809@tcp(127.0.0.1:3306)/rpcsystem")
	//defer utils.DB.Close()
	err := utils.DB.Ping()
	if err != nil {
		fmt.Println("mysql open err:", err.Error())
		return
	}
}
