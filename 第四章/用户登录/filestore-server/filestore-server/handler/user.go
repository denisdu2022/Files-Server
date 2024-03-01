package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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

//SigninHandler: 登录接口

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	//判断请求方式
	//GET请求
	if r.Method == http.MethodGet {
		//读取signin.html文件
		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//返回客户端浏览器HTML文件
		w.Write(data)
		//http.Redirect(w, r, "./static/view/signin.html", http.StatusFound)
		return
	}
	//获取post请求参数
	r.ParseForm()
	//获取用户名和密码
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	fmt.Println(username + password)
	//对传入的密码进行加密,加密和数据库加密的对比
	encpasswd := util.Sha1([]byte(password + pwd_salt))
	fmt.Println(encpasswd)
	//1.检验用户名和密码
	pwdChecked := dblayer.UserSignin(username, encpasswd)
	fmt.Println(pwdChecked)
	//如果校验不通过直接返回
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}

	//2.生成访问凭证(token)
	token := GenToken(username)
	fmt.Println(token)
	upRes := dblayer.UpdateToken(username, token)
	if upRes {
		w.Write([]byte("FAILED"))
		return
	}

	//3.登录成功后重定向到首页
	w.Write([]byte("http://" + r.Host + "/static/view/home.html"))

}

func GenToken(useranme string) string {
	//40位字符串: md5(useranme+timestamp+token_salt)+timestamp[:8]
	//格式化时间戳
	ts := fmt.Sprintf("%x", time.Now().Unix())
	//前32位
	tokenPrefix := util.MD5([]byte(useranme + ts + "_tokensalt"))
	//拼接,返回
	return tokenPrefix + ts[:8]

}
