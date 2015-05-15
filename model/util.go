package model

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	cookies []*http.Cookie
)

// func init() {
// 	cookies = getCookieFromFile("cookies")
// }

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
	for _, v := range cookies {
		req.AddCookie(v)
	}

	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != 200 {
		log.Fatalf("状态码: %d\n获取异常,请休息一会儿再试", resp.StatusCode)
	}
	defer resp.Body.Close()

	if useCookie {
		cookies = resp.Cookies()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

// func setCookieToFile(name string) {
// 	io.buffer
// }

// func getCookieFromFile(name string) []*http.Cookie {
// 	ret, err := ioutil.ReadFile(name)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	content := string(ret)
// 	lines := strings.Split(content, "\n")
// 	for _, line := range lines {

// 	}
// }

func SaveList(name string, songs SongsList) {
	ioutil.WriteFile(name, []byte(songs.String()), 0777)
}
