package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

func main() {
	// Create a Resty Client
	//client := resty.New()
	//
	//type AuthSuccess struct {
	//	Status string `json:"status"`
	//	Msg    string `json:"msg"`
	//}
	//type AuthError struct {
	//}
	//var success AuthSuccess
	//// POST Map, default is JSON content type. No need to set one
	//// Login Request
	//resp, err := client.R().
	//	SetHeaders(map[string]string{"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8"}).
	//	SetBody("action=user_login&username=xxx&password=xxx").
	//	SetResult(&success).    // or SetResult(AuthSuccess{}).
	//	SetError(&AuthError{}). // or SetError(AuthError{}).
	//	Post("https://xueke58.com/wp-admin/admin-ajax.php")
	//if err != nil {
	//	panic(err)
	//}
	//for i, cookie := range resp.Cookies() {
	//	fmt.Printf("cookie%d:%v\n", i, cookie)
	//}

	for i := 0; i < 300; i++ {
		go QianDao()
	}
	// QianDao()
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

	first := "wordpress_020636737623d67f4b561ed84ef6d683=52lhg%7C1712143054%7COCg2aYxJUI68Q9gr9j4xj48aFjTrsdCEjxIMha6SKt2%7C9a4c4f750e76bf50e185f01273d2f1cbe5e01e269f3026cd313f7cbe48ece4ac;"
	second := "PHPSESSID=li5jdi2dkaigi7peqgcsfcpe8f;"
	three := "wordpress_logged_in_020636737623d67f4b561ed84ef6d683=52lhg%7C1712143054%7COCg2aYxJUI68Q9gr9j4xj48aFjTrsdCEjxIMha6SKt2%7Caaed5937d9bcc2a8fb39178b03e6861d99aa314e453c1df72066087963bbaa75"
	ck := first + second + three
	// POST Map, default is JSON content type. No need to set one
	// Login Request
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
		panic(err)
	}
	//fmt.Printf("resp:%v\n", resp)
	fmt.Printf("status:%s,msg:%s\n", success.Status, success.Msg)
}
