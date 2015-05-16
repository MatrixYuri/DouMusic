package model

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

const (
	loginApi      = "http://douban.fm/j/login"
	captchaApi    = "http://douban.fm/j/new_captcha"
	captchaImgApi = "http://douban.fm/misc/captcha?size=l&id=%s" //size=s小号m中号l大号

	likeListApi = "http://douban.fm/j/play_record?ck=%s&spbid=::%s&type=liked&start=%d"
	startFM     = "http://douban.fm?start=%sg%sg"
	playlistApi = "http://douban.fm/j/mine/playlist?type=n&sid=&pt=0.0&channel=0&from=mainsite&kbps=64"
)

var (
	jar       *cookiejar.Jar
	cookieUrl *url.URL
)

func init() {
	cookieUrl, _ = url.ParseRequestURI("http://www.douban.fm")
	jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	getCookieFromFile()
}

func httpDo(url string, args string, method string, useCookie bool) string {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, strings.NewReader(args))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "douban.fm")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36")
	req.Header.Set("Referer", "http://douban.fm/mine?type=liked")

	for _, v := range jar.Cookies(cookieUrl) {
		req.AddCookie(v)
	}

	resp, err := client.Do(req)

	if err != nil {
		return ""
	}
	if resp.StatusCode != 200 {
		log.Printf("请求[%s]时遇到异常\n", url)
		log.Fatalf("状态码: \033[0;31m%d\033[0m,获取异常,请休息一会儿再试", resp.StatusCode)
	}
	defer resp.Body.Close()

	if useCookie {
		jar.SetCookies(cookieUrl, resp.Cookies())
		setCookieToFile()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func setCookieToFile() {
	lines := ""
	for _, cookie := range jar.Cookies(cookieUrl) {
		lines = cookie.Name + "=" + cookie.Value + "\n"
	}
	ioutil.WriteFile("cookie.tmp", []byte(lines), 0664)
}

func getCookieFromFile() {
	ret, err := ioutil.ReadFile("cookie.tmp")
	if err != nil {
		return
	}
	content := string(ret)
	lines := strings.Split(content, "\n")
	var cookies []*http.Cookie
	for _, line := range lines {
		word := strings.Split(line, "=")
		if len(word) <= 1 {
			continue
		}
		cookies = append(cookies, &http.Cookie{Name: word[0], Value: word[1]})
	}
	jar.SetCookies(cookieUrl, cookies)
}

func SaveList(name string, songs SongsList) {
	ioutil.WriteFile(name, []byte(songs.String()), 0664)
}
