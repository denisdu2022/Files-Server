package backup_db

import (
	"filestore-server/db/mysql"
	"fmt"
	"log"
)

//OnFileUploadFinished: 文件上传后保存文件元信息到数据库

func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {

	//使用Prepare 可以防止SQL注入
	stmt, err := mysql.DBInit().Prepare("insert ignore into tbl_file(file_sha1,file_name,file_size,file_addr,status) values (?,?,?,?,1)")

	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed to prepare statement,err: " + err.Error())
		return false
	}

	defer stmt.Close()

	//执行sql语句
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {

		fmt.Println(err.Error())
		return false
	}

	//判断file_sha1文件重复插入的情况
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false
}
