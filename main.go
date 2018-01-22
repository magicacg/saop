package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/headzoo/surf"
)

func regexp0(s1, s2 string) []string {
	r1 := regexp.MustCompile(s2)
	r2 := r1.FindAllString(s1, -1)
	return r2
}
func regexp1(s1, s2, s3 string) string {
	r1 := regexp.MustCompile(s2)
	r2 := r1.ReplaceAllString(s1, s3)
	return r2
}

func main() {
	var html, pixiv_id, password string
	var p, date, date1, date2 string
	var p1 int

	fi, err := os.Open("setting.txt") //读取文件夹里的setting.txt
	if err != nil {
		log.Println("No file: setting.txt!")
		os.Exit(0)
	}
	defer fi.Close()              //关闭文件
	fd, err := ioutil.ReadAll(fi) //IO
	txtbody := string(fd)
	pixiv_id = regexp1(regexp0(txtbody, `"pixiv_id"=[^\n]+`)[0], `("pixiv_id"=|\n|\f|\r)`, "")
	password = regexp1(regexp0(txtbody, `"password"=[^\n]+`)[0], `("password"=|\n|\f|\r)`, "")

	p = regexp1(regexp0(txtbody, `"p"=(\w+|)`)[0], `("p"=|\n|\f|\r)`, ``)
	date = regexp1(regexp0(txtbody, `"date"=(\w+|)`)[0], `("date"=|\n|\f|\r)`, "")
	bow := surf.NewBrowser()
	bow.AddRequestHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
	bow.Open("https://accounts.pixiv.net/login?lang=zh")
	post_key := regexp1(regexp0(bow.Body(), `name="post_key" value="\w+`)[0], `name="post_key" value="`, ``)
	if pixiv_id == "" || password == "" {
		pixiv_id = "" //The is sao gong zhu ,Please do not change !
		password = ""   //上面是我用谷歌翻译的，这里的意思是这个内置帐号是骚公主的，求你不要改改她的密码
	}
	data := make(url.Values) //Getform提交的表单
	data["pixiv_id"] = []string{pixiv_id}
	data["password"] = []string{password}
	data["captcha"] = []string{""}
	data["g_recaptcha_response"] = []string{""}
	data["post_key"] = []string{post_key}
	data["source"] = []string{"pc"}
	data["ref"] = []string{"wwwtop_accounts_index"}
	data["return_to"] = []string{"https://www.pixiv.net/"}
	bow.PostForm("https://accounts.pixiv.net/api/login?lang=zh", data)
	if p == "" {
		p1 = 10
	} else {
		p1, _ = strconv.Atoi(p) //字符串转换整形       p1为帖子的ID
	}
	if date == "" {
		for {
			err = bow.Open("https://www.pixiv.net/ranking.php?mode=daily&p=1")
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			} else {
				break
			}
		}
		date1 = regexp1(regexp0(bow.Body(), `pixiv.context.date = "\w+";`)[0], `[^\d]+`, ``) //找到帖子的最新时间
		date2 = "https://www.pixiv.net/ranking.php?mode=daily&p="
	} else {
		date1 = date
		date2 = "https://www.pixiv.net/ranking.php?mode=daily&date=" + date1 + "&p=" //把setting.txt里的时间赋值

	} //date1为时间
	for i := 1; i <= p1; i++ {
		JSON := make(url.Values) //Getform提交的表单
		JSON["mode"] = []string{"daily"}
		JSON["p"] = []string{strconv.Itoa(i)}
		JSON["format"] = []string{"json"}
		for {
			err := bow.OpenForm(date2+strconv.Itoa(i), JSON) //打开链接
			if err != nil {
				log.Println("Get url err!")
				time.Sleep(2 * time.Second)
				continue
			} else {
				break
			}
		}
		illust := regexp0(bow.Body(), `illust_id&#34;:\w+`)

		t := len(illust)
		for j := 0; j < t; j++ {
			illust_id := regexp1(illust[j], `illust_id&#34;:`, "")

			for {
				err = bow.Open("https://www.pixiv.net/member_illust.php?mode=medium&illust_id=" + illust_id)
				if err != nil {
					continue
				} else {
					break
				}
			}

			urlimage := regexp0(bow.Body(), `(class="original-image"|class="read-more js-click-trackable"|<div class="player toggle">)`)
			t1 := len(urlimage)
			for v := 0; v < t1; v++ {

				if urlimage[v] == `class="original-image"` {
					urlimage1 := regexp1(regexp0(bow.Body(), `data-src="[^\"]+" class="original-image"`)[0], `(data-src=|"|class="original-image|\s| )`, ``)
					fmt.Println(urlimage1)
					html += urlimage1 + "\n"
					break
				} else if urlimage[v] == `class="read-more js-click-trackable"` {
					urlimage1 := regexp0(bow.Body(), `href="[^\"]+" target="[^\"]+" rel="[^\"]+" class="read-more js-click-trackable"`)
					urlimage2 := regexp1(regexp0(urlimage1[0], `href="[^\"]+"`)[0], `(href=|"|amp;)`, ``)
					for {
						err = bow.Open("https://www.pixiv.net" + urlimage2)
						if err != nil {
							time.Sleep(2 * time.Second)
							continue
						} else {
							break
						}
					}
					urlimage3 := regexp0(bow.Body(), `(<a href="[^\"]+" target="[^\"]+" class="full-size-container _ui-tooltip"|class="panel-container visible")`)
					if urlimage3[0] == `class="panel-container visible"` {
						continue
					}
					urlimage4 := regexp1(regexp0(urlimage3[0], `href="[^\"]+"`)[0], `(href=|"|amp;)`, ``)
					for {
						err = bow.Open("https://www.pixiv.net" + urlimage4)
						if err != nil {
							time.Sleep(2 * time.Second)
							continue
						} else {
							break
						}
					}
					urlimage5 := regexp1(regexp0(bow.Body(), `<img src="[^\"]+"`)[0], `(<img src=|")`, ``)
					fmt.Println(urlimage5)
					html += urlimage5 + "\n"
					break
				} else {
					break
				}
			}
		}
	}
	f0, _ := os.Create(date1 + ".txt") //写入到date1时间的TXT文件里
	f0.WriteString(html)               //写入内容
	defer f0.Close()                   //关闭文件
	log.Println(date1 + ".txt	-OK!")
	os.Exit(0)
}
