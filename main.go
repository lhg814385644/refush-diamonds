package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"refush-diamonds/config"
	"sync"
)

func main() {
	// TODO: 读取配置文件的目录
	config.ParseConfig("./bin/")
	if config.C == nil {
		panic("Unmarshal Config  Failed")
	}
	cookie := Login(config.C)
	var wg sync.WaitGroup
	loopNum := 100
	if config.C.ConcurrentNum > 0 {
		loopNum = config.C.ConcurrentNum
	}
	for i := 0; i < loopNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			QianDao(cookie)
		}()
	}
	wg.Wait()
}

// Cookie 登录成功后返回的cookie
type Cookie struct {
	first  string
	second string
	third  string
}

func Login(cfg *config.Config) Cookie {
	// Create a Resty Client
	client := resty.New()

	type AuthSuccess struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
	type AuthError struct {
	}
	var success AuthSuccess
	// POST Map, default is JSON content type. No need to set one
	body := fmt.Sprintf("action=user_login&username=%s&password=%s", cfg.UserName, cfg.Password)
	// Login Request
	resp, err := client.R().
		SetHeaders(map[string]string{"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8"}).
		SetBody(body).
		SetResult(&success).    // or SetResult(AuthSuccess{}).
		SetError(&AuthError{}). // or SetError(AuthError{}).
		Post("https://xueke58.com/wp-admin/admin-ajax.php")
	if err != nil {
		panic(err)
	}
	// for i, cookie := range resp.Cookies() {
	//	fmt.Printf("cookie%d:%v\n", i, cookie)
	// }
	first := resp.Cookies()[2].Name + "=" + resp.Cookies()[2].Value + ";"
	second := resp.Cookies()[0].Name + "=" + resp.Cookies()[0].Value + ";"
	third := resp.Cookies()[3].Name + "=" + resp.Cookies()[3].Value
	return Cookie{
		first:  first,
		second: second,
		third:  third,
	}
}

// QianDao 签到
func QianDao(cookie Cookie) {
	client := resty.New()

	type AuthSuccess struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
	type AuthError struct {
	}
	var success AuthSuccess

	first := cookie.first
	second := cookie.second
	three := cookie.third
	ck := first + second + three

	// POST Map, default is JSON content type. No need to set one
	_, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
			"Cookie":       ck,
		}).
		SetBody("action=user_qiandao").
		SetResult(&success).    // or SetResult(AuthSuccess{}).
		SetError(&AuthError{}). // or SetError(AuthError{}).
		Post("https://xueke58.com/wp-admin/admin-ajax.php")
	if err != nil {
		// TODO: 不要panic,该错误直接不处理
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("status:%s,msg:%s\n", success.Status, success.Msg)
}
