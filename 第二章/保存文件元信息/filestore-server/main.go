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

	//启动监听
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Printf("Failed to start server,err:%s", err.Error())
		return
	}

}
