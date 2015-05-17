package source

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	songSearchApi_163 = "http://music.163.com/api/search/get/web"
	songDetailApi_163 = "http://music.163.com/api/song/detail"
)

type Netease struct {
	Song
}

func (s *Netease) Search(key string) (Song, bool) {
	arg := fmt.Sprintf("type=1&s=%s&offset=0&limit=30", key)

	client := &http.Client{}
	req, err := http.NewRequest("POST", songSearchApi_163, strings.NewReader(arg))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://music.163.com/")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return Song{}, false
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	type searchRes struct {
		Code    int
		Message string
		Result  struct {
			Songs []Netease
		}
	}
	var res searchRes
	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&res)
	if res.Code != 200 {
		fmt.Println(res.Message)
		return Song{}, false
	}

	return Netease.Song, true
}

func (s *Netease) GetDatil(song *Song) bool {
	return true
}
