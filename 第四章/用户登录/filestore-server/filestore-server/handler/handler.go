package handler

import (
	"encoding/json"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// UploadHandler: 处理文件上传

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	//判断请求方式
	if r.Method == "GET" {
		//返回上传HTML页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		//接收文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data,err:%s\n", err.Error())
			return
		}
		defer file.Close()

		//文件元信息
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			//Linux Location
			//Location: "/tmp/" + head.Filename,
			//Win Location
			Location: "D:/tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		//创建文件句柄
		//newFile, err := os.Create("/tmp/" + head.Filename)
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file,err:%s\n", err.Error())
			return
		}
		defer newFile.Close()

		//将接收到的文件拷贝到存储的目录下
		//_, err = io.Copy(newFile, file)
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file,err:%s\n", err.Error())
			return
		}

		//获取FileSha1
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)

		//更新文件元信息
		//meta.UploadFileMeta(fileMeta)
		_ = meta.UpdateFileMetaDB(fileMeta)

		//重定向,到上传成功页面
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}

}

//UploadSucHandler: 上传已完成

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

//GetFileMetaHandler: 获取文件源信息

func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic("err: " + err.Error())
		return
	}

	//获取传入的hash
	fileHash := r.Form["filehash"][0]

	//fMeta := meta.GetFileMeta(fileHash)
	fMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//序列化
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(fMeta.Location, fMeta.FileSha1)
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	//解析参数
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	//打开文件
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	//将文件加载到内存
	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//在浏览器中下载文件,需要设置head
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

//FileMetaUpdateHandler:更新元信息接口(重命名)

func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//使用文件sha1获取文件对象
	curFileMeta := meta.GetFileMeta(fileSha1)
	//文件重命名
	curFileMeta.FileName = newFileName
	//更新文件
	meta.UploadFileMeta(curFileMeta)

	//序列化
	jData, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//文件重命名后返回正确的状态码
	w.WriteHeader(http.StatusOK)

	w.Write(jData)

}

//FileDeleteHandler  删除文件

func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileSha1 := r.Form.Get("filehash")

	fMeta := meta.GetFileMeta(fileSha1)
	//删除文件系统中的文件(物理删除)
	_ = os.Remove(fMeta.Location)

	//删除内存中的文件索引数据
	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
}
