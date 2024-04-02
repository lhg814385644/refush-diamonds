package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"sync"
)

func main() {
	Login()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			QianDao()
		}()
	}
	wg.Wait()
	//QianDao()
}

func Login() {
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
	// Login Request
	resp, err := client.R().
		SetHeaders(map[string]string{"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8"}).
		SetBody("action=user_login&username=toki_l&password=toki_l$").
		SetResult(&success). // or SetResult(AuthSuccess{}).
		SetError(&AuthError{}). // or SetError(AuthError{}).
		Post("https://xueke58.com/wp-admin/admin-ajax.php")
	if err != nil {
		panic(err)
	}
	for i, cookie := range resp.Cookies() {
		fmt.Printf("cookie%d:%v\n", i, cookie)
	}
}

// QianDao 签到
func QianDao() {
	client := resty.New()

	type AuthSuccess struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
	type AuthError struct {
	}
	var success AuthSuccess

	first := "wordpress_020636737623d67f4b561ed84ef6d683=toki_l%7C1712192552%7CYnxaQeA8mgR34L36cnYfvtMpaSqbyfR5y46aLySfsu2%7Ce61c53f88286c6226d8dd46f0c4f50bd02ba582f7388e3e988db0abf83bd489e;"
	second := "PHPSESSID=xxxxasasz;Hm_lvt_f45a3f7cb8733f05fe13ce759d760423=1712018688;wordpress_test_cookie=WP%20Cookie%20check;" // TODO:如何获取？
	three := "wordpress_logged_in_020636737623d67f4b561ed84ef6d683=toki_l%7C1712192552%7CYnxaQeA8mgR34L36cnYfvtMpaSqbyfR5y46aLySfsu2%7C7130e7396cd2c08e8746bd465aa14154cfc227c47d3e11057a85f3ff415cf120; Hm_lpvt_f45a3f7cb8733f05fe13ce759d760423=1712022067"
	ck := first + second + three

	// POST Map, default is JSON content type. No need to set one
	_, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
			"Cookie":       ck,
		}).
		SetBody("action=user_qiandao").
		SetResult(&success). // or SetResult(AuthSuccess{}).
		SetError(&AuthError{}). // or SetError(AuthError{}).
		Post("https://xueke58.com/wp-admin/admin-ajax.php")
	if err != nil {
		// TODO: 不要panic,该错误直接不处理
		fmt.Printf("err:%v\n", err)
		return
	}
	// fmt.Printf("resp:%v\n", resp)
	// TODO: Why????
	fmt.Printf("status:%s,msg:%s\n", success.Status, success.Msg)
}
