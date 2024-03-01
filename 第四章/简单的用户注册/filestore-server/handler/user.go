package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"io/ioutil"
	"net/http"
)

//盐值

const (
	pwd_salt = "*#8901"
)

//SignupHandler : 处理用户注册请求

func SignupHandler(w http.ResponseWriter, r *http.Request) {

	//判断请求方式
	//GET请求
	if r.Method == http.MethodGet {
		//读取signup.html文件
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//返回客户端浏览器HTML文件
		w.Write(data)
		return
	}
	//POST请求
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")
	//判断用户名密码的校验
	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("Invalid parameter!"))
		return
	}

	//对用户名和密码进行加密
	enc_passwd := util.Sha1([]byte(passwd + pwd_salt))
	suc := dblayer.UserSignup(username, enc_passwd)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("Failed"))
	}
}
