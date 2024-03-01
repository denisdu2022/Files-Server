package backup_db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
	"log"
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

//UserSignin : 判断用户名密码是否一致

func UserSignin(useranme string, encpwd string) bool {
	//查询语句
	stmt, err := mydb.GetDB().Prepare("select * from tbl_user where user_name =?  limit 1")
	if err != nil {
		log.Fatal(err)
		return false
	}
	//执行查询
	rows, err := stmt.Query(useranme)
	if err != nil {
		log.Fatal(err)
		return false
	} else if rows == nil {
		fmt.Println("username not found: " + useranme)
		return false
	}
	//对比密码
	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

//UpdateToken : 更新用户登录的token

func UpdateToken(username string, token string) bool {
	//SQL语句 ,replace覆盖插入,因为是user_name是唯一主键,以插入的最新的user_token为凭证
	stmt, err := mydb.GetDB().Prepare("replace into tbl_user_token(user_name,user_token) values (?,?)")
	if err != nil {
		log.Fatal(err)
		return false
	}

	defer stmt.Close()
	//执行查询
	_, err = stmt.Exec(username, token)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

//查询用户信息

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

func GetUserInfo(username string) (User, error) {
	//实例化
	user := User{}
	//数据库查询
	stmt, err := mydb.GetDB().Prepare("select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
		//如果出错了,返回空的User和err
		return user, err
	}
	defer stmt.Close()
	//执行sql语句
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		log.Fatal(err)
		return user, err
	}
	return user, nil
}
