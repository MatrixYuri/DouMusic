package model

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	cookies []*http.Cookie
)

func httpDo(url string, args string, method string) string {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, strings.NewReader(args))
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "douban.fm")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36")
	req.Header.Set("Referer", "http://douban.fm/mine?type=liked")
	for _, v := range cookies {
		req.AddCookie(v)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	if url == loginApi {
		cookies = resp.Cookies()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	return string(body)
}

func SaveList(name string, songs SongsList) {
	ioutil.WriteFile(name+"喜欢的歌曲.txt", []byte(songs.String()), 0777)
}
