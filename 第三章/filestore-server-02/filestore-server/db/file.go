package backup_db

import (
	"database/sql"
	mydb "filestore-server/db/mysql"
	"fmt"
	"log"
)

//OnFileUploadFinished: 文件上传后保存文件元信息到数据库

func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {

	//使用Prepare 可以防止SQL注入
	stmt, err := mydb.GetDB().Prepare("insert ignore into tbl_file(file_sha1,file_name,file_size,file_addr,status) values (?,?,?,?,1)")

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

//查询文件元信息

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

//GetFileMeta : 从mysql获取文件元信息

func GetFileMeta(filehash string) (*TableFile, error) {

	//数据库查询
	stmt, err := mydb.GetDB().Prepare("SELECT file_sha1,file_addr,file_name,file_size FROM tbl_file WHERE file_sha1 = ? AND status =1 LIMIT 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	//执行查询返回单条记录
	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &tfile, err

}
