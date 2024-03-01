package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

// 获取DB对象
var db *sql.DB

//初始化数据库连接

func initDB() *sql.DB {

	//数据库连接
	db, _ := sql.Open("mysql", "root:xxxxxxx@tcp(127.0.0.1:3306)/fileserver?charset=utf8")

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

func GetDB() *sql.DB {
	db := initDB()
	return db

}

func ParseRows(rows *sql.Rows) []map[string]interface{} {

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		//检测错误
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
