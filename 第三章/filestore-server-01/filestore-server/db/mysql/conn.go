package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

// 获取DB对象
var db *sql.DB

//初始化数据库连接

func DBInit() *sql.DB {

	//数据库连接
	db, _ := sql.Open("mysql", "root:xxxxxx@tcp(127.0.0.1:3306)/fileserver?charset=utf8")

	//设置最大连接数
	db.SetMaxOpenConns(500)
	//连接测试
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql,err: " + err.Error())
		//强制退出进程
		os.Exit(1)
	}

	return db
}
