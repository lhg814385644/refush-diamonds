package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"refush-diamonds/config"
	myZap "refush-diamonds/zap"
	"strings"
	"sync"
)

func main() {
	_ = myZap.SetLevelFromString("info")
	myZap.InitZap("dev")

	// TODO: 读取配置文件的目录
	logger := zap.L().With(zap.String("service", "钻石"))
	config.ParseConfig("./bin/")
	logger.Info("服务启动。。。")
	if config.C == nil {
		panic("Unmarshal Config  Failed")
	}
	logger.Info("服务启动成功....")

	cookie := Login(config.C)
	// ShowDiamonds(cookie)
	var wg sync.WaitGroup
	runTimes := 100
	if config.C.ConcurrentNum > 0 {
		runTimes = config.C.ConcurrentNum
	}

	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			QianDao(cookie)
		}()
	}
	wg.Wait()

	/**********************************************/
	// 创建通道跟踪成功请求
	//successCh := make(chan bool, 10)
	//successCount := 0
	//p, err := ants.NewPool(20)
	//if err != nil {
	//	fmt.Printf("NewPool err:%v\n", err)
	//}
	//// 监控成功通道以达到阈值
	//go func() {
	//	for {
	//		select {
	//		case result := <-successCh:
	//			if result {
	//				successCount++
	//			} else {
	//				successCount = 0
	//			}
	//			if successCount >= 5 {
	//				p.Release() // 关闭池并释放worker队列
	//				fmt.Println("成功阈值达到，停止后续请求")
	//				return
	//			}
	//			// TODO: 需要设一个默认能退出的标志 暂时不管
	//		}
	//	}
	//}()
	//for i := 0; i < runTimes; i++ {
	//	wg.Add(1)
	//	err := p.Submit(func() {
	//		defer wg.Done()              // submit的函数内部必须wg.Done
	//		QianDaoV2(cookie, successCh) // TODO:需要有返回值，并接收
	//	})
	//	if err != nil {
	//		fmt.Printf("submit err:%v\n", err)
	//		return
	//	}
	//}
	//// 确保池能正常释放
	//if !p.IsClosed() {
	//	p.Release()
	//}
	//wg.Wait()
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
		SetResult(&success). // or SetResult(AuthSuccess{}).
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
			"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
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
	fmt.Printf("status:%s,msg:%s\n", success.Status, success.Msg)
}

func QianDaoV2(cookie Cookie, successCh chan bool) {
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
			"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
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
	// TODO:发送信号表明成功请求
	if success.Status == "1" {
		successCh <- true
	} else {
		successCh <- false
	}
	fmt.Printf("status:%s,msg:%s\n", success.Status, success.Msg)
}

// ShowDiamonds 抓取当前余额
func ShowDiamonds(cookie Cookie) {
	client := resty.New()

	type AuthSuccess struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
	type AuthError struct {
	}
	//var success AuthSuccess

	first := cookie.first
	second := cookie.second
	three := cookie.third
	ck := first + second + three

	// POST Map, default is JSON content type. No need to set one
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type": "text/html; charset=UTF-8",
			"Cookie":       ck,
		}).
		// SetBody("action=user_qiandao").
		// SetResult(&success).    // or SetResult(AuthSuccess{}).
		SetError(&AuthError{}). // or SetError(AuthError{}).
		Get("https://xueke58.com/user")
	if err != nil {
		// TODO: 不要panic,该错误直接不处理
		fmt.Printf("err:%v\n", err)
		return
	}
	//fmt.Printf("status:%s,msg:%s\n", success.Status, success.Msg)
	// 解析 HTML 文档

	doc, err := html.Parse(strings.NewReader(string(resp.Body())))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(doc.Attr)

	// 查找 "现有余额" 元素
	//balance := doc.Find(".jinbi")

	// 获取文本内容
	//text := balance.Text()

	// 打印结果
	//fmt.Println("现有余额:", text)

	//fileName := "index.html"
	//
	//bs, err := ioutil.ReadFile(fileName)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//text := string(bs)
	//
	//doc, err := html.Parse(strings.NewReader(text))
	//
	//if err != nil {
	//
	//	log.Fatal(err)
	//}

	var data []string

	doTraverse(doc, &data, "span")
	// fmt.Println(data)
	for i, datum := range data {
		fmt.Printf("span%d=%v\n", i, datum)
	}
}

func doTraverse(doc *html.Node, data *[]string, tag string) {

	var traverse func(n *html.Node, tag string) *html.Node

	traverse = func(n *html.Node, tag string) *html.Node {

		for c := n.FirstChild; c != nil; c = c.NextSibling {

			if c.Type == html.TextNode && c.Parent.Data == tag {

				*data = append(*data, c.Data)
			}

			res := traverse(c, tag)

			if res != nil {

				return res
			}
		}

		return nil
	}

	traverse(doc, tag)
}
