package model

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

type Song struct {
	Id           string `json:"id"`            //歌曲id
	Title        string `json:"title"`         //歌曲名
	Artist       string `json:"artist"`        //歌手
	Liked        bool   `json:"liked"`         //是否红心
	SubjectTitle string `json:"subject_title"` //专辑名
	Picture      string `json:"picture"`       //专辑封面
	Path         string `json:"path"`          //所属专辑地址
	IsDone       bool   `json:"-"`             //是否已下载
}

type SongsList struct {
	Total    int    `json:"total"`
	Start    int    `json:"start"`
	PerPage  int    `json:"per_page"`
	Songs    []Song `json:"songs"`
	SongType string `json:"song_type"`
}

func (song Song) String() string {
	return fmt.Sprintf("id=%s|done=%v\n", song.Id, song.IsDone)
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
	for _, cookie := range cookies {
		if cookie.Name == "bid" {
			bid = cookie.Value
			break
		}
	}
	url := fmt.Sprintf(likeListApi, ck, bid, 0)
	ret := httpDo(url, "", "GET")
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
			ret := httpDo(url, "", "GET")
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

func (list *SongsList) Download() {

}
