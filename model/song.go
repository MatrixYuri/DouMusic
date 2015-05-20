package model

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

type Song struct {
	Id           string `json:"id"`            //歌曲id,也是sid
	Ssid         string `json:"ssid"`          //歌曲ssid,关键
	Title        string `json:"title"`         //歌曲名
	Artist       string `json:"artist"`        //歌手
	Liked        bool   `json:"liked"`         //是否红心
	SubjectTitle string `json:"subject_title"` //专辑名
	Picture      string `json:"picture"`       //专辑封面
	Path         string `json:"path"`          //所属专辑地址
	Url          string `json:"url"`           //下载地址
}

type SongsList struct {
	Total    int    `json:"total"`
	Start    int    `json:"start"`
	PerPage  int    `json:"per_page"`
	Songs    []Song `json:"songs"`
	SongType string `json:"song_type"`
}

func (song Song) String() string {
	return fmt.Sprintf("%s|%s\n", song.Title, song.Url)
}

func (list SongsList) String() string {
	ret := ""
	for _, v := range list.Songs {
		ret += v.String()
	}
	return ret
}

//获取全部红心列表
//ck 为User.Ck
func GetList(ck string) SongsList {
	var bid string
	for _, cookie := range jar.Cookies(cookieUrl) {
		if cookie.Name == "bid" {
			bid = cookie.Value
			break
		}
	}
	url := fmt.Sprintf(likeListApi, ck, bid, 0)
	ret := httpDo(url, "", "GET", false)
	var songs SongsList
	err := json.NewDecoder(strings.NewReader(ret)).Decode(&songs)
	if err != nil {
		fmt.Println(err)
		return SongsList{}
	}

	pageNum := int(math.Floor(float64(songs.Total) / float64(songs.PerPage)))
	var c chan []Song = make(chan []Song)
	for i := 1; i <= pageNum; i++ {
		go func(index int) {
			url := fmt.Sprintf(likeListApi, ck, bid, index*15)
			ret := httpDo(url, "", "GET", false)
			if ret == "{\"songs\":[]}" {
				return
			}
			var page SongsList
			err := json.NewDecoder(strings.NewReader(ret)).Decode(&page)
			if err != nil {
				return
			}
			c <- page.Songs
		}(i)
	}

	for i := 1; i <= pageNum; i++ {
		list := <-c
		songs.Songs = append(songs.Songs, list...)
	}
	close(c)

	return songs
}

func Download() bool {
	return false
}
