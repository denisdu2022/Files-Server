package main

import (
	"filestore-server/handler"
	"fmt"
	"net/http"
)

func main() {

	//路由
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	http.HandleFunc("/user/signup", handler.SignupHandler)

	//启动监听
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Printf("Failed to start server,err:%s", err.Error())
		return
	}

}
