package backup_db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
)

//UserSignup: 通过用户名和密码注册用户

func UserSignup(username string, passwd string) bool {
	//数据库插入用户名和密码
	stmt, err := mydb.GetDB().Prepare("insert ignore into tbl_user (user_name,user_pwd) values (?,?)")
	if err != nil {
		fmt.Println("Failed to insert,err: " + err.Error())
		return false
	}
	defer stmt.Close()
	//执行SQL语句
	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert,err: " + err.Error())
	}
	//判断是否有重复插入
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}
