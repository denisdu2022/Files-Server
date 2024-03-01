package meta

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

//GetFileMeta:通过fileSha1获取文件元信息的对象接口

func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

//RemoveFileMeta: 通过fileSha1移除元信息

func RemoveFileMeta(filesha1 string) {
	delete(fileMetas, filesha1)
}
