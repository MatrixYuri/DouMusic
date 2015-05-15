package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
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

func (song *Song) GetSsid() {
	doc, err := goquery.NewDocument(song.Path)
	if err != nil {
		log.Fatal(err)
	}
	song.Ssid = doc.Find("li#"+song.Id).AttrOr("data-ssid", "0")
}

var lock sync.Mutex

func (song *Song) GetDownloadLink() bool {
	url := fmt.Sprintf(startFM, song.Id, song.Ssid)

	lock.Lock()
	httpDo(url, "", "GET", true)
	ret := httpDo(playlistApi, "", "GET", false)
	for ret == "" {
		log.Println("\033[0;31m重试:\033[0m\t" + song.Title)
		ret = httpDo(playlistApi, "", "GET", false)
	}
	lock.Unlock()

	type playlist struct {
		R     int    `json:"r"`
		Songs []Song `json:"song"`
	}
	var temp playlist
	err := json.NewDecoder(strings.NewReader(ret)).Decode(&temp)
	if err != nil || temp.R != 0 {
		log.Println(err)
		return false
	}
	for _, v := range temp.Songs {
		if v.Ssid == song.Ssid {
			song.Url = v.Url
			log.Println("得到:\t" + song.Title)
			return true
		}
	}
	log.Println("\033[0;31m未匹配:\033[0m\t" + song.Title)
	return false
}

func (list *SongsList) Download() {
	var getWG sync.WaitGroup
	for i, _ := range list.Songs {
		getWG.Add(1)
		go func(index int) {
			list.Songs[index].GetSsid()
			list.Songs[index].GetDownloadLink()
			SaveList("downloadTask.txt", *list)
			defer getWG.Done()
		}(i)
	}
	getWG.Wait()

	fmt.Println("\033[0;32m====下载链接整理完毕，开始下载====\033[0m")

	var downWG sync.WaitGroup
	for _, song := range list.Songs {
		fileName := "./download/" + song.Title + ".mp3"
		if _, err := os.Stat(fileName); os.IsNotExist(err) && song.Url != "" {
			downWG.Add(1)
			go func(name string, url string) {
				data := httpDo(url, "", "GET", false)

				file, _ := os.Create(name)
				defer file.Close()
				io.Copy(file, bytes.NewReader([]byte(data)))
				log.Printf("保存:\t\033[0;33m%s\033[0m\n", name)
				defer downWG.Done()
			}(fileName, song.Url)
		} else {
			log.Printf("跳过:\t\033[0;37m%s\033[0m\n", song.Title)
		}
	}
	downWG.Wait()
}
