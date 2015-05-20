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
	searchApi_163 = "http://music.163.com/api/search/get/web"
	detailApi_163 = "http://music.163.com/api/song/detail"
)

type Netease struct {
	Song    `json:"-"`
	Id      string `json:"id"`
	Name    string `json:"name"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Album struct {
		Name string `json:"name"`
	} `json:"album"`
	Duration int `json:"duration"`
}

func (s *Netease) Search(key string) (Song, bool) {
	arg := fmt.Sprintf("type=1&s=%s&offset=0&limit=30", key)

	client := &http.Client{}
	req, err := http.NewRequest("POST", searchApi_163, strings.NewReader(arg))
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
		Code    int `json:"code"`
		Message string
		Result  struct {
			SongCount int       `json:"songCount"`
			Songs     []Netease `json:"songs"`
		} `json:"result"`
	}
	var res searchRes
	err = json.Unmarshal(body, &res)
	log.Println(res.Result)
	if res.Code != 200 {
		fmt.Println(res.Message)
		return Song{}, false
	}

	// for _, each := range res.Result.Songs {
	//     if each. {

	//     }
	// }

	return Song{}, true
}

func (s *Netease) GetDatil(song *Song) bool {
	return true
}
