package meta

import (
	mydb "filestore-server/db"
	"fmt"
)

//FileMeta:文件元信息结构

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

//初始化文件元信息

func init() {
	fileMetas = make(map[string]FileMeta)
}

//UploadFileMeta:新增/更新文件元信息

func UploadFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

//UpdateFileMetaDB: 新增/更新文件元信息到mysql中

func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(
		fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

//GetFileMeta:通过fileSha1获取文件元信息的对象接口

func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

//从mysql返回文件元信息

func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	//调用数据库查询方法
	tfile, err := mydb.GetFileMeta(fileSha1)
	if err != nil {
		fmt.Println(err.Error())
		return FileMeta{}, err
	}
	//转换结构
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta, nil
}

//RemoveFileMeta: 通过fileSha1移除元信息

func RemoveFileMeta(filesha1 string) {
	delete(fileMetas, filesha1)
}
