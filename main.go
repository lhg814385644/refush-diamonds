package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

func main() {
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
		SetBody("action=user_login&username=xxx&password=xxx").
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
